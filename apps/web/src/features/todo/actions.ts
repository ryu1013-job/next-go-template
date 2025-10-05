'use server'

import type { ServerActionState } from '~/lib/create-server-action'
import { createTodo } from '~/generated/client'
import { createServerAction } from '~/lib/create-server-action'
import { TodoSchema } from './schema'

export async function addTodo(_prevState: ServerActionState, formData: FormData): Promise<ServerActionState> {
  return createServerAction(
    formData,
    ['title', 'description'],
    TodoSchema,
    async (data) => {
      await createTodo(data)
    },
    'TODOの作成に失敗しました',
  )
}
