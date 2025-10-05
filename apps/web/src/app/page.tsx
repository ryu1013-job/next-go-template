import { AddForm } from '~/features/todo/add-form'
import { listTodos } from '~/generated/client'

export default async function Home() {
  const todos = await listTodos()

  return (
    <div className="max-w-2xl mx-auto p-4">
      <AddForm />
      <pre>{JSON.stringify(todos, null, 2)}</pre>
    </div>
  )
}
