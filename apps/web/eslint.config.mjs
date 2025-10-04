import antfu from '@antfu/eslint-config'

export default antfu({
  ignores: ['**/generated/**', '**/node_modules/**', '**/.next/**'],
})
