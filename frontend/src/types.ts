export interface Todo {
  id: number
  title: string
  completed: boolean
  createdAt: string
  updatedAt: string
}

export interface TodoPatch {
  title?: string
  completed?: boolean
}
