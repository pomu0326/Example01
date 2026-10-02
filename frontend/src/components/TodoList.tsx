import type { Todo, TodoPatch } from '../types'
import { TodoItem } from './TodoItem'

interface Props {
  todos: Todo[]
  onUpdate: (id: number, patch: TodoPatch) => Promise<boolean>
  onDelete: (id: number) => Promise<boolean>
}

export function TodoList({ todos, onUpdate, onDelete }: Props) {
  if (todos.length === 0) {
    return <p className="empty">タスクはありません</p>
  }
  return (
    <ul className="todo-list">
      {todos.map((todo) => (
        <TodoItem key={todo.id} todo={todo} onUpdate={onUpdate} onDelete={onDelete} />
      ))}
    </ul>
  )
}
