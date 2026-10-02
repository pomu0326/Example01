import { cleanup, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import * as api from './api'
import type { Todo } from './types'

vi.mock('./api')

function todo(overrides: Partial<Todo>): Todo {
  return {
    id: 1,
    title: 'task',
    completed: false,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    ...overrides,
  }
}

describe('App', () => {
  beforeEach(() => {
    vi.mocked(api.listTodos).mockResolvedValue([
      todo({ id: 2, title: '牛乳を買う' }),
      todo({ id: 1, title: 'メールを返す', completed: true }),
    ])
  })

  afterEach(() => {
    cleanup()
    vi.resetAllMocks()
  })

  it('shows todos and the remaining count', async () => {
    render(<App />)
    expect(await screen.findByText('牛乳を買う')).toBeInTheDocument()
    expect(screen.getByText('メールを返す')).toBeInTheDocument()
    expect(screen.getByText('残り 1 件')).toBeInTheDocument()
  })

  it('shows an empty message when there are no todos', async () => {
    vi.mocked(api.listTodos).mockResolvedValue([])
    render(<App />)
    expect(await screen.findByText('タスクはありません')).toBeInTheDocument()
  })

  it('adds a todo to the top of the list', async () => {
    vi.mocked(api.createTodo).mockResolvedValue(todo({ id: 3, title: '本を読む' }))
    const user = userEvent.setup()
    render(<App />)
    await screen.findByText('牛乳を買う')

    await user.type(screen.getByLabelText('新しいタスク'), '  本を読む  ')
    await user.click(screen.getByRole('button', { name: '追加' }))

    expect(api.createTodo).toHaveBeenCalledWith('本を読む')
    const items = await screen.findAllByRole('listitem')
    expect(items[0]).toHaveTextContent('本を読む')
    expect(screen.getByLabelText('新しいタスク')).toHaveValue('')
  })

  it('keeps the input and shows the error when adding fails', async () => {
    vi.mocked(api.createTodo).mockRejectedValue(new Error('title must be 1 to 200 characters'))
    const user = userEvent.setup()
    render(<App />)
    await screen.findByText('牛乳を買う')

    await user.type(screen.getByLabelText('新しいタスク'), 'x')
    await user.click(screen.getByRole('button', { name: '追加' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('title must be 1 to 200 characters')
    expect(screen.getByLabelText('新しいタスク')).toHaveValue('x')
  })

  it('toggles completion', async () => {
    vi.mocked(api.updateTodo).mockResolvedValue(todo({ id: 2, title: '牛乳を買う', completed: true }))
    const user = userEvent.setup()
    render(<App />)

    await user.click(await screen.findByLabelText('牛乳を買う を完了にする'))

    expect(api.updateTodo).toHaveBeenCalledWith(2, { completed: true })
    expect(await screen.findByText('残り 0 件')).toBeInTheDocument()
  })

  it('edits a title', async () => {
    vi.mocked(api.updateTodo).mockResolvedValue(todo({ id: 2, title: '豆乳を買う' }))
    const user = userEvent.setup()
    render(<App />)
    const item = (await screen.findByText('牛乳を買う')).closest('li')!

    await user.click(within(item).getByRole('button', { name: '編集' }))
    const input = within(item).getByLabelText('タイトルを編集')
    await user.clear(input)
    await user.type(input, '豆乳を買う{Enter}')

    expect(api.updateTodo).toHaveBeenCalledWith(2, { title: '豆乳を買う' })
    expect(await screen.findByText('豆乳を買う')).toBeInTheDocument()
  })

  it('deletes a todo', async () => {
    vi.mocked(api.deleteTodo).mockResolvedValue(undefined)
    const user = userEvent.setup()
    render(<App />)
    const item = (await screen.findByText('牛乳を買う')).closest('li')!

    await user.click(within(item).getByRole('button', { name: '削除' }))

    expect(api.deleteTodo).toHaveBeenCalledWith(2)
    expect(screen.queryByText('牛乳を買う')).not.toBeInTheDocument()
  })

  it('shows an error when loading fails', async () => {
    vi.mocked(api.listTodos).mockRejectedValue(new Error('HTTP 500'))
    render(<App />)
    expect(await screen.findByRole('alert')).toHaveTextContent('HTTP 500')
  })
})
