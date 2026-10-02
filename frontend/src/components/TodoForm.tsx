import { useState, type FormEvent } from 'react'

interface Props {
  onAdd: (title: string) => Promise<boolean>
}

export function TodoForm({ onAdd }: Props) {
  const [title, setTitle] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    const trimmed = title.trim()
    if (!trimmed) return
    setSubmitting(true)
    const ok = await onAdd(trimmed)
    setSubmitting(false)
    if (ok) setTitle('')
  }

  return (
    <form onSubmit={handleSubmit} className="todo-form">
      <input
        aria-label="新しいタスク"
        placeholder="やることを入力"
        value={title}
        maxLength={200}
        onChange={(e) => setTitle(e.target.value)}
      />
      <button type="submit" disabled={submitting || !title.trim()}>
        追加
      </button>
    </form>
  )
}
