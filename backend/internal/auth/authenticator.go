package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	jwksURL  string
	issuer   string
	audience string

	httpClient *http.Client

	mu   sync.RWMutex
	keys map[string]*ecdsa.PublicKey
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewAuthenticator(jwksURL, issuer, audience string) (*Authenticator, error) {
	a := &Authenticator{
		jwksURL:  jwksURL,
		issuer:   issuer,
		audience: audience,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		keys: make(map[string]*ecdsa.PublicKey),
	}

	if err := a.refreshKeys(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Authenticator) UserIDFromRequest(r *http.Request) (string, error) {
	tokenString, err := bearerTokenFromRequest(r)
	if err != nil {
		return "", err
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		a.keyFunc,
		jwt.WithValidMethods([]string{"ES256"}),
		jwt.WithIssuer(a.issuer),
		jwt.WithAudience(a.audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.Subject == "" {
		return "", errors.New("sub claim is empty")
	}

	return claims.Subject, nil
}

func bearerTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with Bearer")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}

func (a *Authenticator) keyFunc(token *jwt.Token) (interface{}, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, errors.New("kid header is missing")
	}

	key := a.getKey(kid)
	if key != nil {
		return key, nil
	}

	// kid が見つからない場合、鍵ローテーションの可能性があるので再取得
	if err := a.refreshKeys(); err != nil {
		return nil, err
	}

	key = a.getKey(kid)
	if key == nil {
		return nil, fmt.Errorf("public key not found for kid: %s", kid)
	}

	return key, nil
}

func (a *Authenticator) getKey(kid string) *ecdsa.PublicKey {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.keys[kid]
}

func (a *Authenticator) refreshKeys() error {
	res, err := a.httpClient.Get(a.jwksURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch jwks: status=%d", res.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return err
	}

	keys := make(map[string]*ecdsa.PublicKey)

	for _, jwk := range jwks.Keys {
		if jwk.Kid == "" {
			continue
		}

		// Supabase ECC(P-256) / ES256 前提
		if jwk.Kty != "EC" || jwk.Crv != "P-256" || jwk.Alg != "ES256" {
			continue
		}

		key, err := ecdsaPublicKeyFromJWK(jwk)
		if err != nil {
			return err
		}

		keys[jwk.Kid] = key
	}

	if len(keys) == 0 {
		return errors.New("no valid ES256 public keys found in jwks")
	}

	a.mu.Lock()
	a.keys = keys
	a.mu.Unlock()

	return nil
}

func ecdsaPublicKeyFromJWK(jwk JWK) (*ecdsa.PublicKey, error) {
	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, err
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, err
	}

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	curve := elliptic.P256()

	if !curve.IsOnCurve(x, y) {
		return nil, errors.New("public key is not on P-256 curve")
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}
