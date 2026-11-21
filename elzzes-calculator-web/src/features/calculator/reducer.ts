import {
  backspace,
  displayEntry,
  entryValue,
  inputDigit,
  inputDot,
  toggleSign
} from './entry'
import {
  binarySymbols,
  operandText,
  operandValue,
  unaryLabel
} from './expression'
import { formatNumber } from './format'
import {
  initialState,
  type Action,
  type ActiveRequest,
  type BinaryOperation,
  type CalculatorKey,
  type CalculatorState,
  type ClientError,
  type Operand,
  type RequestPlan,
  type UnaryOperation,
  type Value
} from './types'

export function reducer(
  state: CalculatorState,
  action: Action
): CalculatorState {
  switch (action.type) {
    case 'press':
      return press(state, action.key)
    case 'requestSucceeded':
      return succeed(state, action.id, action.result)
    case 'requestFailed':
      return fail(state, action.id, action.error)
    default:
      return state
  }
}

function press(state: CalculatorState, key: CalculatorKey): CalculatorState {
  if (key === 'clear') {
    return { ...initialState, sequence: state.sequence }
  }
  if (state.phase === 'calculating') {
    return state
  }
  if (key === 'equals') {
    return pressEquals(state)
  }
  if (key === 'sqrt' || key === 'percent') {
    return pressUnary(state, key)
  }
  if (
    key === 'add' ||
    key === 'subtract' ||
    key === 'multiply' ||
    key === 'divide' ||
    key === 'power'
  ) {
    return pressBinary(state, key)
  }
  if (state.phase === 'error') {
    if (key === 'backspace' || key === 'sign') {
      return state
    }
    return startFreshEntry(state, key)
  }
  if (state.phase === 'result') {
    if (key === 'backspace' || key === 'sign') {
      return state
    }
    return startFreshEntry(state, key)
  }
  if (key === 'backspace') {
    return editEntry(state, backspace)
  }
  if (key === 'sign') {
    return editEntry(state, toggleSign)
  }
  if (key === '.') {
    return pressDot(state)
  }
  return pressDigit(state, key)
}

function startFreshEntry(
  state: CalculatorState,
  key: CalculatorKey
): CalculatorState {
  const text =
    key === '.'
      ? inputDot('')
      : key === '0' || /^[1-9]$/.test(key)
        ? inputDigit('', key)
        : ''
  return {
    ...initialState,
    phase: 'idle',
    sequence: state.sequence,
    current: { kind: 'entry', text }
  }
}

function pressDigit(state: CalculatorState, digit: string): CalculatorState {
  if (state.current.kind === 'value') {
    return startSecondEntry(state, digit)
  }
  return editEntry(state, (text) => inputDigit(text, digit))
}

function pressDot(state: CalculatorState): CalculatorState {
  if (state.current.kind === 'value') {
    return startSecondEntry(state, '.')
  }
  return editEntry(state, inputDot)
}

function startSecondEntry(
  state: CalculatorState,
  key: string
): CalculatorState {
  if (!state.pending || !state.first) {
    return {
      ...initialState,
      phase: 'idle',
      sequence: state.sequence,
      current: {
        kind: 'entry',
        text: key === '.' ? inputDot('') : inputDigit('', key)
      }
    }
  }

  const text = key === '.' ? inputDot('') : inputDigit('', key)
  return {
    ...state,
    phase: 'idle',
    current: { kind: 'entry', text },
    second: null,
    secondLabel: text,
    error: null,
    retry: null,
    expression: `${state.first.text} ${binarySymbols[state.pending]}`
  }
}

function editEntry(
  state: CalculatorState,
  edit: (text: string) => string
): CalculatorState {
  if (state.current.kind !== 'entry') {
    return state
  }

  const text = edit(state.current.text)
  const hasSecond = text !== ''
  const secondLabel = hasSecond ? text : null
  const expression =
    state.pending && state.first
      ? `${state.first.text} ${binarySymbols[state.pending]}`
      : state.expression

  return {
    ...state,
    current: { kind: 'entry', text },
    second: null,
    secondLabel,
    expression
  }
}

function pressBinary(
  state: CalculatorState,
  operation: BinaryOperation
): CalculatorState {
  if (state.phase === 'error') {
    if (state.error?.kind === 'domain') {
      const frozen = operandText(state.current)
      const first: Value = { value: Number(frozen), text: frozen }
      return {
        ...state,
        phase: 'idle',
        current: { kind: 'value', ...first },
        first,
        second: null,
        secondLabel: null,
        pending: operation,
        error: null,
        retry: null,
        expression: `${frozen} ${binarySymbols[operation]}`
      }
    }

    const first = state.first ?? valueFromOperand(state.current)
    return {
      ...state,
      phase: 'idle',
      current: { kind: 'value', ...first },
      first,
      second: null,
      secondLabel: null,
      pending: operation,
      error: null,
      retry: null,
      expression: `${first.text} ${binarySymbols[operation]}`
    }
  }

  if (state.phase === 'result') {
    const first = valueFromOperand(state.current)
    return {
      ...state,
      phase: 'idle',
      current: { kind: 'value', ...first },
      first,
      pending: operation,
      second: null,
      secondLabel: null,
      expression: `${first.text} ${binarySymbols[operation]}`
    }
  }

  if (!state.pending) {
    const first = valueFromOperand(state.current)
    return {
      ...state,
      phase: 'idle',
      current: { kind: 'value', ...first },
      first,
      pending: operation,
      second: null,
      secondLabel: null,
      expression: `${first.text} ${binarySymbols[operation]}`
    }
  }

  if (!hasSecond(state)) {
    const first = state.first ?? valueFromOperand(state.current)
    return {
      ...state,
      phase: 'idle',
      current: { kind: 'value', ...first },
      first,
      pending: operation,
      second: null,
      secondLabel: null,
      expression: `${first.text} ${binarySymbols[operation]}`
    }
  }

  const first = state.first as Value
  const secondOperand = state.second ?? state.current
  const secondText = state.secondLabel ?? operandText(secondOperand)
  const expression = `${first.text} ${binarySymbols[state.pending]} ${secondText}`
  return startRequest(state, {
    operation: state.pending,
    operands: [first.value, operandValue(secondOperand)],
    continuation: { kind: 'chain', next: operation },
    expression,
    display: secondText
  })
}

function pressEquals(state: CalculatorState): CalculatorState {
  if (state.phase === 'result') {
    return state
  }

  if (state.phase === 'error') {
    if (state.error?.kind === 'domain' || !state.retry) {
      return state
    }
    return startRequest(state, state.retry)
  }

  if (!state.pending || !hasSecond(state)) {
    return state
  }

  const first = state.first as Value
  const secondOperand = state.second ?? state.current
  const secondText = state.secondLabel ?? operandText(secondOperand)
  return startRequest(state, {
    operation: state.pending,
    operands: [first.value, operandValue(secondOperand)],
    continuation: { kind: 'equals' },
    expression: `${first.text} ${binarySymbols[state.pending]} ${secondText}`,
    display: secondText
  })
}

function pressUnary(
  state: CalculatorState,
  operation: UnaryOperation
): CalculatorState {
  if (state.phase === 'error') {
    return state
  }

  if (state.pending) {
    if (!hasSecond(state)) {
      return state
    }
    const first = state.first as Value
    const source = state.second ?? state.current
    const sourceText = state.secondLabel ?? operandText(source)
    const label = unaryLabel(operation, sourceText)
    return startRequest(state, {
      operation,
      operands: [operandValue(source)],
      continuation: { kind: 'unary-second', label },
      expression: `${first.text} ${binarySymbols[state.pending]} ${label}`,
      display: sourceText
    })
  }

  const sourceText = operandText(state.current)
  return startRequest(state, {
    operation,
    operands: [operandValue(state.current)],
    continuation: { kind: 'unary-current' },
    expression: unaryLabel(operation, sourceText),
    display: sourceText
  })
}

function startRequest(
  state: CalculatorState,
  plan: RequestPlan
): CalculatorState {
  const request: ActiveRequest = { ...plan, id: state.sequence + 1 }
  return {
    ...state,
    phase: 'calculating',
    request,
    retry: null,
    error: null,
    sequence: request.id
  }
}

function succeed(
  state: CalculatorState,
  id: number,
  result: number
): CalculatorState {
  if (state.phase !== 'calculating' || state.request?.id !== id) {
    return state
  }

  const request = state.request
  const text = formatNumber(result)
  const value: Value = { value: Number(text), text }
  const continuation = request.continuation

  if (continuation.kind === 'chain') {
    return {
      ...state,
      phase: 'idle',
      current: { kind: 'value', ...value },
      first: value,
      second: null,
      secondLabel: null,
      pending: continuation.next,
      expression: `${text} ${binarySymbols[continuation.next]}`,
      error: null,
      request: null,
      retry: null
    }
  }

  if (continuation.kind === 'unary-second') {
    return {
      ...state,
      phase: 'idle',
      current: { kind: 'value', ...value },
      second: { kind: 'value', ...value },
      secondLabel: continuation.label,
      expression: request.expression,
      error: null,
      request: null,
      retry: null
    }
  }

  return {
    ...state,
    phase: 'result',
    current: { kind: 'value', ...value },
    first: value,
    second: null,
    secondLabel: null,
    pending: null,
    expression: `${request.expression} = ${text}`,
    error: null,
    request: null,
    retry: null
  }
}

function fail(
  state: CalculatorState,
  id: number,
  error: ClientError
): CalculatorState {
  if (state.phase !== 'calculating' || state.request?.id !== id) {
    return state
  }

  const request = state.request

  if (error.kind === 'domain') {
    const frozen = request.display
    return {
      ...state,
      phase: 'error',
      current: { kind: 'value', value: Number(frozen), text: frozen },
      first: null,
      second: null,
      secondLabel: null,
      pending: null,
      expression: request.expression,
      error,
      request: null,
      retry: null
    }
  }

  return {
    ...state,
    phase: 'error',
    expression: request.expression,
    error,
    request: null,
    retry: {
      operation: request.operation,
      operands: request.operands,
      continuation: request.continuation,
      expression: request.expression,
      display: request.display
    }
  }
}

function hasSecond(state: CalculatorState): boolean {
  return (
    state.second !== null ||
    (state.current.kind === 'entry' && state.current.text !== '')
  )
}

function valueFromOperand(operand: Operand): Value {
  if (operand.kind === 'value') {
    return { value: operand.value, text: operand.text }
  }
  return {
    value: entryValue(operand.text),
    text: displayEntry(operand.text)
  }
}
