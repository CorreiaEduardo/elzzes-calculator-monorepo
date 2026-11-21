import { errorForCode } from '../copy'
import type { ClientError, ServerErrorCode } from '../types'

const domain422 = new Set<ServerErrorCode>([
  'DIVISION_BY_ZERO',
  'NEGATIVE_SQUARE_ROOT',
  'NON_FINITE_RESULT'
])

export function serverError(
  status: number,
  code: ServerErrorCode
): ClientError {
  if (status === 422 && domain422.has(code)) {
    return errorForCode(code)
  }
  if (status >= 500) {
    return errorForCode('INTERNAL_ERROR')
  }
  if (status >= 400 && status < 500) {
    return errorForCode('INVALID_CALCULATION')
  }
  return errorForCode('INVALID_RESPONSE')
}

export function transportError(
  code: 'NETWORK_ERROR' | 'TIMEOUT' | 'INVALID_RESPONSE'
): ClientError {
  return errorForCode(code)
}
