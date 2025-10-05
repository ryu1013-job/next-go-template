import type { TodoCreateInput } from '~/generated/client'
import { minLength, object, optional, pipe, string, transform } from 'valibot'

export const TodoSchema = object({
  title: pipe(
    string(),
    transform(input => input.trim()),
    minLength(1, 'タイトルを入力してください'),
  ),
  description: optional(string()),
  dueDate: optional(string()),
} satisfies Record<keyof TodoCreateInput, any>)

export type TodoInput = TodoCreateInput
