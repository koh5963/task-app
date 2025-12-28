FROM golang:1.25.5-alpine3.21

WORKDIR /app

# COPY go.mod go.sum ./ This project still has no go.sum.
COPY go.mod ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .
EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
