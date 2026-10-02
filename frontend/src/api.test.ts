import { afterEach, describe, expect, it, vi } from 'vitest'
import { createTodo, deleteTodo, listTodos } from './api'

function mockFetch(status: number, body?: unknown) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(body === undefined ? null : JSON.stringify(body), { status }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('api', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('lists todos', async () => {
    const fetchMock = mockFetch(200, [{ id: 1, title: 'a' }])
    await expect(listTodos()).resolves.toEqual([{ id: 1, title: 'a' }])
    expect(fetchMock).toHaveBeenCalledWith('/api/todos', expect.anything())
  })

  it('sends a JSON body when creating', async () => {
    const fetchMock = mockFetch(201, { id: 1, title: 'a' })
    await createTodo('a')
    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('/api/todos')
    expect(init.method).toBe('POST')
    expect(init.headers).toEqual({ 'Content-Type': 'application/json' })
    expect(JSON.parse(init.body)).toEqual({ title: 'a' })
  })

  it('handles 204 No Content', async () => {
    mockFetch(204)
    await expect(deleteTodo(1)).resolves.toBeUndefined()
  })

  it('throws the server error message', async () => {
    mockFetch(404, { error: 'todo not found' })
    await expect(deleteTodo(1)).rejects.toThrow('todo not found')
  })

  it('falls back to the status when the error body is not JSON', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('oops', { status: 502 })))
    await expect(listTodos()).rejects.toThrow('HTTP 502')
  })
})
