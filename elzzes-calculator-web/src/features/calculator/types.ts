export type Phase = 'idle' | 'calculating' | 'result' | 'error'

export type BinaryOperation =
  'add' | 'subtract' | 'multiply' | 'divide' | 'power'
export type UnaryOperation = 'sqrt' | 'percent'
export type Operation = BinaryOperation | UnaryOperation

export const serverErrorCodes = [
  'INVALID_JSON',
  'UNKNOWN_FIELD',
  'UNKNOWN_OPERATION',
  'INVALID_OPERAND_COUNT',
  'INVALID_OPERAND',
  'PAYLOAD_TOO_LARGE',
  'INVALID_CONTENT_TYPE',
  'DIVISION_BY_ZERO',
  'NEGATIVE_SQUARE_ROOT',
  'NON_FINITE_RESULT',
  'NOT_FOUND',
  'METHOD_NOT_ALLOWED',
  'INTERNAL_ERROR'
] as const

export type ServerErrorCode = (typeof serverErrorCodes)[number]

export type ClientErrorCode =
  | ServerErrorCode
  | 'INVALID_CALCULATION'
  | 'NETWORK_ERROR'
  | 'TIMEOUT'
  | 'INVALID_RESPONSE'

export interface ClientError {
  kind: 'domain' | 'transient'
  code: ClientErrorCode
  copy: string
}

export type Operand =
  | { kind: 'entry'; text: string }
  | { kind: 'value'; value: number; text: string }

export type Value = { value: number; text: string }

export type CalculatorKey =
  | '0'
  | '1'
  | '2'
  | '3'
  | '4'
  | '5'
  | '6'
  | '7'
  | '8'
  | '9'
  | '.'
  | 'backspace'
  | 'sign'
  | 'clear'
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percent'
  | 'equals'

export type Continuation =
  | { kind: 'equals' }
  | { kind: 'chain'; next: BinaryOperation }
  | { kind: 'unary-current' }
  | { kind: 'unary-second'; label: string }

export interface RequestPlan {
  operation: Operation
  operands: [number] | [number, number]
  continuation: Continuation
  expression: string
  display: string
}

export interface ActiveRequest extends RequestPlan {
  id: number
}

export interface CalculatorState {
  phase: Phase
  current: Operand
  first: Value | null
  second: Operand | null
  secondLabel: string | null
  pending: BinaryOperation | null
  expression: string
  error: ClientError | null
  request: ActiveRequest | null
  retry: RequestPlan | null
  sequence: number
}

export type Action =
  | { type: 'press'; key: CalculatorKey }
  | { type: 'requestSucceeded'; id: number; result: number }
  | { type: 'requestFailed'; id: number; error: ClientError }

export const initialState: CalculatorState = {
  phase: 'idle',
  current: { kind: 'entry', text: '' },
  first: null,
  second: null,
  secondLabel: null,
  pending: null,
  expression: '',
  error: null,
  request: null,
  retry: null,
  sequence: 0
}
