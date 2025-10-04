import { defineConfig } from 'orval'

export default defineConfig({
  todoApi: {
    input: '../../packages/schema/openapi.yaml',
    output: {
      mode: 'single',
      target: 'src/generated/client.ts',
      client: 'fetch',
      clean: true,
      override: {
        mutator: {
          path: 'src/lib/fetcher.ts',
          name: 'apiFetch',
        },
      },
    },
  },
})
