// TaskList.tsx
import { useState } from 'react'
import type { Task } from '../models/Task'
import './App.css'
import { supabase } from '../../../shared/SupabaseClient'

type TaskStatus = 'TODO' | 'DOING' | 'DONE'

function TaskList() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [status, setStatus] = useState<TaskStatus>('TODO')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const getAuthHeaders = async () => {
    const { data } = await supabase.auth.getSession()
    const token = data.session?.access_token

    if (!token) {
      throw new Error('not logged in')
    }

    return {
      Authorization: `Bearer ${token}`,
    }
  }

  const handleLoadTasks = async () => {
    setLoading(true)
    setError(null)

    try {
      const headers = await getAuthHeaders()

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

  const handleCreateTask = async () => {
    setLoading(true)
    setError(null)

    try {
      const headers = await getAuthHeaders()

      const res = await fetch('http://localhost:8080/tasks', {
        method: 'POST',
        headers: {
          ...headers,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title,
          description,
          status,
        }),
      })

      if (!res.ok) {
        const body = await res.text()
        console.log('status:', res.status)
        console.log('body:', body)
        throw new Error(`failed to create task: ${res.status}`)
      }

      const newTasks: Task[] = await res.json()
      setTasks(newTasks)

      setTitle('')
      setDescription('')
      setStatus('TODO')
    } catch (err) {
      setError('タスク作成に失敗しました')
      console.log(err)
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateStatus = async (taskId: string, nextStatus: TaskStatus) => {
    setLoading(true)
    setError(null)

    try {
      const headers = await getAuthHeaders()

      const res = await fetch(`http://localhost:8080/tasks/${taskId}/${nextStatus}`, {
        method: 'PATCH',
        headers,
      })

      if (!res.ok) {
        const body = await res.text()
        console.log('status:', res.status)
        console.log('body:', body)
        throw new Error(`failed to update task: ${res.status}`)
      }

      await handleLoadTasks()
    } catch (err) {
      setError('タスク更新に失敗しました')
      console.log(err)
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteTask = async (taskId: string) => {
    setLoading(true)
    setError(null)

    try {
      const headers = await getAuthHeaders()

      const res = await fetch(`http://localhost:8080/tasks/${taskId}`, {
        method: 'DELETE',
        headers,
      })

      if (!res.ok) {
        const body = await res.text()
        console.log('status:', res.status)
        console.log('body:', body)
        throw new Error(`failed to delete task: ${res.status}`)
      }

      await handleLoadTasks()
    } catch (err) {
      setError('タスク削除に失敗しました')
      console.log(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <section className="task-form">
        <h2>タスク作成</h2>

        <div className="task-form-row">
          <label htmlFor="task-title">タイトル</label>
          <input
            id="task-title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="タスク名"
          />
        </div>

        <div className="task-form-row">
          <label htmlFor="task-description">説明</label>
          <input
            id="task-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="説明"
          />
        </div>

        <div className="task-form-row">
          <label htmlFor="task-status">ステータス</label>
          <select
            id="task-status"
            value={status}
            onChange={(e) => setStatus(e.target.value as TaskStatus)}
          >
            <option value="TODO">TODO</option>
            <option value="DOING">DOING</option>
            <option value="DONE">DONE</option>
          </select>
        </div>

        <button onClick={handleCreateTask} disabled={loading || title.trim() === ''}>
          作成
        </button>
      </section>

      <section className="task-list-section">
        <button onClick={handleLoadTasks} disabled={loading}>
          {loading ? '処理中...' : '一覧取得'}
        </button>

        {error && <p className="error-message">{error}</p>}

        <ul>
          {tasks.map((task) => (
            <li key={task.id}>
              <div className="task-header">
                <strong>{task.title}</strong>
                <span className="task-status">{task.status}</span>
              </div>

              <p>{task.description}</p>

              <div className="task-actions">
                <button onClick={() => handleUpdateStatus(task.id, 'TODO')} disabled={loading}>
                  TODO
                </button>
                <button onClick={() => handleUpdateStatus(task.id, 'DOING')} disabled={loading}>
                  DOING
                </button>
                <button onClick={() => handleUpdateStatus(task.id, 'DONE')} disabled={loading}>
                  DONE
                </button>
                <button onClick={() => handleDeleteTask(task.id)} disabled={loading}>
                  削除
                </button>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </>
  )
}

export default TaskList