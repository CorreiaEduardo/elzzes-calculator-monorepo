import type { ClientError, Operation } from '../types'
import { calculateErrorSchema, calculateResponseSchema } from './schemas'
import { serverError, transportError } from './error'

const CALCULATE_PATH = '/api/v1/calculate'
const TIMEOUT_MS = 10_000

export type CalculateResult =
  | { kind: 'success'; result: number }
  | { kind: 'error'; error: ClientError }
  | { kind: 'aborted' }

export async function calculate(
  operation: Operation,
  operands: number[],
  signal?: AbortSignal
): Promise<CalculateResult> {
  if (signal?.aborted) {
    return { kind: 'aborted' }
  }

  const controller = new AbortController()
  let timedOut = false

  const timeout = setTimeout(() => {
    timedOut = true
    controller.abort()
  }, TIMEOUT_MS)

  const onAbort = () => controller.abort()
  signal?.addEventListener('abort', onAbort, { once: true })

  try {
    const response = await fetch(
      new URL(CALCULATE_PATH, window.location.origin),
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ operation, operands }),
        signal: controller.signal
      }
    )

    let body: unknown
    try {
      body = await response.json()
    } catch {
      return { kind: 'error', error: transportError('INVALID_RESPONSE') }
    }

    if (response.ok) {
      const parsed = calculateResponseSchema.safeParse(body)
      if (!parsed.success) {
        return { kind: 'error', error: transportError('INVALID_RESPONSE') }
      }
      return { kind: 'success', result: parsed.data.result }
    }

    const parsed = calculateErrorSchema.safeParse(body)
    if (!parsed.success) {
      return { kind: 'error', error: transportError('INVALID_RESPONSE') }
    }
    return {
      kind: 'error',
      error: serverError(response.status, parsed.data.error.code)
    }
  } catch {
    if (timedOut) {
      return { kind: 'error', error: transportError('TIMEOUT') }
    }
    if (signal?.aborted) {
      return { kind: 'aborted' }
    }
    return { kind: 'error', error: transportError('NETWORK_ERROR') }
  } finally {
    clearTimeout(timeout)
    signal?.removeEventListener('abort', onAbort)
  }
}
