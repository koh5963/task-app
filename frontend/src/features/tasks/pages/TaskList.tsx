import { useState } from 'react'
import type { Task } from '../models/Task'
import './App.css'

function TaskList() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleLoadTasks = async () => {
    setLoading(true)
    setError(null)

    try {
      const res = await fetch('http://localhost:8080/tasks')

      if (!res.ok) {
        throw new Error('failed to load tasks')
      }

      const data: Task[] = await res.json()
      setTasks(data)
    } catch (err) {
      setError('タスク一覧の取得に失敗しました')
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
