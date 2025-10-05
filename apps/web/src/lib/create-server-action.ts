import type { BaseIssue, BaseSchema, SafeParseResult } from 'valibot'
import { safeParse as valibotSafeParse } from 'valibot'

export interface ServerActionState {
  errors?: Record<string, string>
  ok?: boolean
}

/**
 * 汎用的なServer Action関数
 * @param formData フォームデータ
 * @param keys 取得するフォームフィールドのキー
 * @param schema Valibotスキーマ
 * @param action 実行するアクション（バリデーション成功時）
 * @param errorMessage エラー時のメッセージ（デフォルト: "処理中にエラーが発生しました"）
 */
export async function createServerAction<T extends Record<string, any>, R = any>(
  formData: FormData,
  keys: (keyof T)[],
  schema: BaseSchema<any, any, any>,
  action: (data: T) => Promise<R>,
  errorMessage: string = '処理中にエラーが発生しました',
): Promise<ServerActionState> {
  try {
    const input = getFormDataValues<T>(formData, keys)
    const parsed = valibotSafeParse(schema, input)

    if (!parsed.success) {
      const errors = formatValidationErrors(parsed)
      return { errors }
    }

    await action(parsed.output)
    return { ok: true }
  }
  catch (err) {
    console.error(err)
    return {
      errors: {
        form: errorMessage,
      },
    }
  }
}

/**
 * Valibotのバリデーション結果からフォーム用のエラーオブジェクトを生成する
 */
export function formatValidationErrors(result: SafeParseResult<any>): Record<string, string> {
  if (result.success) {
    return {}
  }

  return Object.fromEntries(
    result.issues.map((issue: BaseIssue<any>) => {
      const path = issue.path && issue.path.length > 0
        ? issue.path.map((p: any) => String(p.key)).join('.')
        : 'form'
      return [path, issue.message]
    }),
  )
}

/**
 * FormDataから指定されたキーの値を安全に取得する
 */
export function getFormDataValue(formData: FormData, key: string): string {
  const value = formData.get(key)
  return typeof value === 'string' ? value : ''
}

/**
 * FormDataから複数のキーの値をオブジェクトとして取得する
 */
export function getFormDataValues<T extends Record<string, any>>(
  formData: FormData,
  keys: (keyof T)[],
): Record<keyof T, string> {
  return Object.fromEntries(
    keys.map(key => [key, getFormDataValue(formData, String(key))]),
  ) as Record<keyof T, string>
}
