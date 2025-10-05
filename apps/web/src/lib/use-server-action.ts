'use client'

import type { ServerActionState } from './create-server-action'
import { useActionState, useEffect, useRef } from 'react'

export interface UseServerActionOptions {
  /**
   * 成功時のコールバック
   */
  onSuccess?: () => void
  /**
   * エラー時のコールバック
   */
  onError?: (error: { message: string, errors?: Record<string, string> }) => void
  /**
   * 処理開始時のコールバック
   */
  onPending?: () => void
}

export interface UseServerActionResult {
  /**
   * Formのactionプロパティにそのまま渡せる関数
   */
  action: (formData: FormData) => void
  /**
   * エラー情報
   */
  error: { message: string, errors?: Record<string, string> } | null
  /**
   * 実行中かどうか
   */
  isPending: boolean
  /**
   * Server Actionの戻り値（状態）
   */
  data: ServerActionState | null
}

/**
 * oRPC風のuseServerActionフック
 * Server Actionを簡単に実行し、状態管理とコールバックをサポート
 */
export function useServerAction(
  action: (prevState: ServerActionState, formData: FormData) => Promise<ServerActionState>,
  options: UseServerActionOptions = {},
): UseServerActionResult {
  const { onSuccess, onError, onPending } = options

  const [state, formAction, isPending] = useActionState<ServerActionState, FormData>(action, {})
  const previousStateRef = useRef<ServerActionState>({})

  const error = state?.errors && Object.keys(state.errors).length > 0
    ? {
        message: Object.values(state.errors)[0] || 'エラーが発生しました',
        errors: state.errors,
      }
    : null

  // 状態変化を監視してコールバックを実行
  useEffect(() => {
    // 初回レンダリングや同じ状態の場合はスキップ
    if (state === previousStateRef.current) {
      return
    }

    previousStateRef.current = state

    if (isPending && onPending) {
      onPending()
    }
    else if (state?.ok && onSuccess) {
      onSuccess()
    }
    else if (state?.errors && Object.keys(state.errors).length > 0 && onError) {
      onError({
        message: Object.values(state.errors)[0] || 'エラーが発生しました',
        errors: state.errors,
      })
    }
  }, [state, isPending, onSuccess, onError, onPending])

  return {
    action: formAction,
    error,
    isPending,
    data: state,
  }
}
