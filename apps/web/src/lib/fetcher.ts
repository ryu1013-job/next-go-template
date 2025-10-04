export type RequestOptions = Omit<RequestInit, 'body'> & { body?: any }

export async function apiFetch<TResponse>(
  url: string,
  options: RequestOptions = {},
): Promise<TResponse> {
  // eslint-disable-next-line node/prefer-global/process
  const base = process.env.NEXT_PUBLIC_API_URL!
  const res = await fetch(base + url, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
    credentials: 'include',
    ...options,
    body: options.body && typeof options.body !== 'string'
      ? JSON.stringify(options.body)
      : options.body,
  })

  if (!res.ok) {
    const text = await res.text().catch(() => '')
    throw new Error(`API ${res.status}: ${text}`)
  }

  return (await res.json()) as TResponse
}
