import { useState, type FormEvent } from 'react'
import type { Todo, TodoPatch } from '../types'

interface Props {
  todo: Todo
  onUpdate: (id: number, patch: TodoPatch) => Promise<boolean>
  onDelete: (id: number) => Promise<boolean>
}

export function TodoItem({ todo, onUpdate, onDelete }: Props) {
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(todo.title)

  async function handleSave(e: FormEvent) {
    e.preventDefault()
    const trimmed = draft.trim()
    if (trimmed && trimmed !== todo.title && !(await onUpdate(todo.id, { title: trimmed }))) {
      return // keep the editor open so the user can retry
    }
    setEditing(false)
  }

  function startEditing() {
    setDraft(todo.title)
    setEditing(true)
  }

  return (
    <li className={todo.completed ? 'todo-item completed' : 'todo-item'}>
      <input
        type="checkbox"
        aria-label={`${todo.title} を完了にする`}
        checked={todo.completed}
        onChange={() => onUpdate(todo.id, { completed: !todo.completed })}
      />
      {editing ? (
        <form onSubmit={handleSave} className="edit-form">
          <input
            aria-label="タイトルを編集"
            value={draft}
            maxLength={200}
            autoFocus
            onChange={(e) => setDraft(e.target.value)}
            onKeyDown={(e) => e.key === 'Escape' && setEditing(false)}
          />
          <button type="submit">保存</button>
        </form>
      ) : (
        <span className="title" onDoubleClick={startEditing}>
          {todo.title}
        </span>
      )}
      {!editing && (
        <button type="button" onClick={startEditing}>
          編集
        </button>
      )}
      <button type="button" onClick={() => onDelete(todo.id)}>
        削除
      </button>
    </li>
  )
}
