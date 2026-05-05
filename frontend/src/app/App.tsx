import { useState, useEffect } from 'react'
import './App.css'
import LoginPage from '../features/auth/pages/LoginPage'
import { supabase } from '../shared/SupabaseClient'
import type { Session } from '@supabase/supabase-js'
import TaskList from '../features/tasks/pages/TaskList'

function App() {
  const [session, setSession] = useState<Session | null>(null)

  useEffect(() => {
    supabase.auth.getSession().then(({ data }) => {
      setSession(data.session)
    })

    const { data: sub } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session)
    })

    return () => sub.subscription.unsubscribe()
  }, [])

  if (!session) return <LoginPage />

  return (
    <>
      <div>
        <p>ログイン済み: {session.user.email}</p>
        <button onClick={() => supabase.auth.signOut()}>logout</button>

        <TaskList />
      </div>
    </>
  )
}


export default App
