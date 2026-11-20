import type { ClientError, ClientErrorCode } from './types'

export const busyText = 'Calculating…'

const domainCodes = new Set<ClientErrorCode>([
  'DIVISION_BY_ZERO',
  'NEGATIVE_SQUARE_ROOT',
  'NON_FINITE_RESULT',
  'INVALID_JSON',
  'UNKNOWN_FIELD',
  'UNKNOWN_OPERATION',
  'INVALID_OPERAND_COUNT',
  'INVALID_OPERAND',
  'PAYLOAD_TOO_LARGE',
  'INVALID_CONTENT_TYPE',
  'NOT_FOUND',
  'METHOD_NOT_ALLOWED',
  'INVALID_CALCULATION'
])

const copies: Record<ClientErrorCode, string> = {
  DIVISION_BY_ZERO: 'Cannot divide by zero.',
  NEGATIVE_SQUARE_ROOT: 'Square root needs a number ≥ 0.',
  NON_FINITE_RESULT: 'That result is too large or undefined.',
  INVALID_JSON: "That calculation isn't valid.",
  UNKNOWN_FIELD: "That calculation isn't valid.",
  UNKNOWN_OPERATION: "That calculation isn't valid.",
  INVALID_OPERAND_COUNT: "That calculation isn't valid.",
  INVALID_OPERAND: "That calculation isn't valid.",
  PAYLOAD_TOO_LARGE: "That calculation isn't valid.",
  INVALID_CONTENT_TYPE: "That calculation isn't valid.",
  NOT_FOUND: "That calculation isn't valid.",
  METHOD_NOT_ALLOWED: "That calculation isn't valid.",
  INVALID_CALCULATION: "That calculation isn't valid.",
  NETWORK_ERROR: "Can't reach the calculator. Press = to retry.",
  TIMEOUT: 'The calculator took too long. Press = to retry.',
  INTERNAL_ERROR: 'Something went wrong. Press = to retry.',
  INVALID_RESPONSE:
    'The calculator returned an unexpected response. Press = to retry.'
}

export function errorForCode(code: ClientErrorCode): ClientError {
  return {
    kind: domainCodes.has(code) ? 'domain' : 'transient',
    code,
    copy: copies[code]
  }
}
