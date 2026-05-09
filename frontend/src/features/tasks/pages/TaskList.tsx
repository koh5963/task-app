import { useState } from 'react'
import type { Task } from '../models/Task'
import './App.css'
import { supabase } from '../../../shared/SupabaseClient'

function TaskList() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleLoadTasks = async () => {
    setLoading(true)
    setError(null)

    try {
      const { data } = await supabase.auth.getSession()
      const token = data.session?.access_token
      const headers = {
        Authorization: `Bearer ${token}`,
      }

      const res = await fetch('http://localhost:8080/tasks', { headers })

      if (!res.ok) {
        const body = await res.text()
        console.log('status:', res.status)
        console.log('body:', body)
        throw new Error(`failed to load tasks: ${res.status}`)
      }

      const tasks: Task[] = await res.json()
      setTasks(tasks)
    } catch (err) {
      setError('タスク一覧の取得に失敗しました')
      console.log(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <button onClick={handleLoadTasks} disabled={loading}>
        {loading ? '取得中...' : '一覧取得'}
      </button>

      {error && <p>{error}</p>}

      <ul>
        {tasks.map((task) => (
          <li key={task.id}>
            <strong>{task.title}</strong>
            <p>{task.description}</p>
            <p>{task.status}</p>
          </li>
        ))}
      </ul>
    </>
  )
}

export default TaskList
