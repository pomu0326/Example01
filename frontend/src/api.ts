import type { Todo, TodoPatch } from './types'

const BASE = '/api/todos'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: init?.body ? { 'Content-Type': 'application/json' } : undefined,
  })
  if (!res.ok) {
    let message = `HTTP ${res.status}`
    try {
      const body = (await res.json()) as { error?: string }
      if (body.error) message = body.error
    } catch {
      // Non-JSON error body; keep the status message.
    }
    throw new Error(message)
  }
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

export function listTodos(): Promise<Todo[]> {
  return request<Todo[]>(BASE)
}

export function createTodo(title: string): Promise<Todo> {
  return request<Todo>(BASE, { method: 'POST', body: JSON.stringify({ title }) })
}

export function updateTodo(id: number, patch: TodoPatch): Promise<Todo> {
  return request<Todo>(`${BASE}/${id}`, { method: 'PATCH', body: JSON.stringify(patch) })
}

export function deleteTodo(id: number): Promise<void> {
  return request<void>(`${BASE}/${id}`, { method: 'DELETE' })
}
