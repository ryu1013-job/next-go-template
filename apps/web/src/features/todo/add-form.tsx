'use client'

import { Button, Form, Input } from '@heroui/react'
import { addToast } from '@heroui/toast'
import { useServerAction } from '~/lib/use-server-action'
import { addTodo } from './actions'

export function AddForm() {
  const { action, isPending, error } = useServerAction(addTodo, {
    onSuccess: () => {
      addToast({
        title: '成功',
        description: 'TODOを作成しました',
        color: 'success',
      })
    },
    onError: (error) => {
      addToast({
        title: 'エラー',
        description: error.message,
        color: 'danger',
      })
    },
  })

  return (
    <Form action={action} validationErrors={error?.errors}>
      <Input
        name="title"
        label="Title"
        isRequired
      />
      <Input name="description" label="Description" />
      <Button type="submit" isLoading={isPending} color="primary">
        Add
      </Button>
    </Form>
  )
}
