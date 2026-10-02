import { useEffect, useState } from 'react'
import * as api from './api'
import { TodoForm } from './components/TodoForm'
import { TodoList } from './components/TodoList'
import type { Todo, TodoPatch } from './types'

function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .listTodos()
      .then(setTodos)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  // Runs an API call, surfacing any failure in the error banner. Resolves to whether it succeeded.
  async function run(action: () => Promise<void>): Promise<boolean> {
    setError(null)
    try {
      await action()
      return true
    } catch (e) {
      setError((e as Error).message)
      return false
    }
  }

  const handleAdd = (title: string) =>
    run(async () => {
      const created = await api.createTodo(title)
      setTodos((prev) => [created, ...prev])
    })

  const handleUpdate = (id: number, patch: TodoPatch) =>
    run(async () => {
      const updated = await api.updateTodo(id, patch)
      setTodos((prev) => prev.map((t) => (t.id === id ? updated : t)))
    })

  const handleDelete = (id: number) =>
    run(async () => {
      await api.deleteTodo(id)
      setTodos((prev) => prev.filter((t) => t.id !== id))
    })

  const remaining = todos.filter((t) => !t.completed).length

  return (
    <main className="app">
      <h1>ToDo</h1>
      <TodoForm onAdd={handleAdd} />
      {error && (
        <p role="alert" className="error">
          {error}
        </p>
      )}
      {loading ? (
        <p>読み込み中...</p>
      ) : (
        <>
          <TodoList todos={todos} onUpdate={handleUpdate} onDelete={handleDelete} />
          <p className="summary">残り {remaining} 件</p>
        </>
      )}
    </main>
  )
}

export default App
