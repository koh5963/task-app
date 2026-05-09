import { useState } from 'react'
import { supabase } from '../../../shared/SupabaseClient'
import './App.css'

function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      const {error} = await supabase.auth.signInWithPassword({ email, password })
      if (error) {
        setError(error.message)
      }
    } catch (err) {
      setError('ログインに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-container">
      <form className="login-form" onSubmit={handleSubmit}>
        <div className="form-row">
          <p className='login-label'>Email</p>
          <input
            className='login-form'
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <div className="form-row">
          <p className='login-label'>Password</p>
          <input
            className='login-form'
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        {error && <p className="error">{error}</p>}

        <button
          type="submit"
          className="button-login"
          disabled={loading}
        >
          {loading ? 'logging in...' : 'login'}
        </button>
      </form>
    </div>
  )
}

export default LoginPage
