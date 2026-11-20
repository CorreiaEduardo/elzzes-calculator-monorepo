import { describe, expect, it } from 'vitest'
import { errorForCode } from './copy'
import { reducer } from './reducer'
import { initialState, type CalculatorKey, type CalculatorState } from './types'

function press(state: CalculatorState, keys: CalculatorKey[]): CalculatorState {
  return keys.reduce(
    (current, key) => reducer(current, { type: 'press', key }),
    state
  )
}

function succeed(state: CalculatorState, result: number): CalculatorState {
  if (!state.request) {
    throw new Error('no request in flight')
  }
  return reducer(state, {
    type: 'requestSucceeded',
    id: state.request.id,
    result
  })
}

function fail(
  state: CalculatorState,
  code: Parameters<typeof errorForCode>[0]
): CalculatorState {
  if (!state.request) {
    throw new Error('no request in flight')
  }
  return reducer(state, {
    type: 'requestFailed',
    id: state.request.id,
    error: errorForCode(code)
  })
}

describe('calculator reducer', () => {
  it('starts at zero with no expression', () => {
    expect(initialState.current).toEqual({ kind: 'entry', text: '' })
    expect(initialState.phase).toBe('idle')
    expect(initialState.expression).toBe('')
  })

  it('applies entry rules', () => {
    let state = press(initialState, ['0', '5', '.', '2'])
    expect(state.current).toEqual({ kind: 'entry', text: '5.2' })

    state = press(state, ['backspace'])
    expect(state.current).toEqual({ kind: 'entry', text: '5.' })

    state = press(state, ['sign'])
    expect(state.current).toEqual({ kind: 'entry', text: '-5.' })

    state = press(state, ['clear'])
    expect(state.current).toEqual({ kind: 'entry', text: '' })
  })

  it('commits the first operand without a request', () => {
    const state = press(initialState, ['1', '2', 'add'])

    expect(state.phase).toBe('idle')
    expect(state.pending).toBe('add')
    expect(state.first).toEqual({ value: 12, text: '12' })
    expect(state.expression).toBe('12 +')
    expect(state.request).toBeNull()
  })

  it('replaces an operator while the second operand is empty', () => {
    const state = press(initialState, ['2', 'add', 'multiply'])

    expect(state.phase).toBe('idle')
    expect(state.pending).toBe('multiply')
    expect(state.expression).toBe('2 ×')
    expect(state.request).toBeNull()
  })

  it('executes equals with exact operands and result text', () => {
    let state = press(initialState, ['1', '2', 'add', '7', 'equals'])

    expect(state.phase).toBe('calculating')
    expect(state.request?.operation).toBe('add')
    expect(state.request?.operands).toEqual([12, 7])
    expect(state.request?.expression).toBe('12 + 7')
    expect(state.request?.continuation).toEqual({ kind: 'equals' })

    state = succeed(state, 19)
    expect(state.phase).toBe('result')
    expect(state.expression).toBe('12 + 7 = 19')
    expect(state.current).toEqual({ kind: 'value', value: 19, text: '19' })
    expect(state.pending).toBeNull()
  })

  it('chains after a completed second operand with one request', () => {
    let state = press(initialState, ['2', 'add', '3', 'add'])

    expect(state.phase).toBe('calculating')
    expect(state.request?.operands).toEqual([2, 3])
    expect(state.request?.continuation).toEqual({
      kind: 'chain',
      next: 'add'
    })

    state = succeed(state, 5)
    expect(state.phase).toBe('idle')
    expect(state.pending).toBe('add')
    expect(state.expression).toBe('5 +')
    expect(state.current.text).toBe('5')
  })

  it('does not execute equals without a complete expression', () => {
    const partial = press(initialState, ['2', 'add'])
    expect(press(partial, ['equals'])).toBe(partial)

    const result = succeed(press(initialState, ['2', 'add', '3', 'equals']), 5)
    expect(press(result, ['equals'])).toBe(result)
  })

  it('applies unary operations to entries and results', () => {
    let state = press(initialState, ['9', 'sqrt'])
    expect(state.request?.operation).toBe('sqrt')
    expect(state.request?.operands).toEqual([9])

    state = succeed(state, 3)
    expect(state.phase).toBe('result')
    expect(state.expression).toBe('√(9) = 3')

    state = press(state, ['percent'])
    expect(state.request?.operation).toBe('percent')
    expect(state.request?.operands).toEqual([3])
    state = succeed(state, 0.03)
    expect(state.expression).toBe('3% = 0.03')
  })

  it('does not apply unary operations to an empty second operand', () => {
    const state = press(initialState, ['1', '2', 'add', 'sqrt'])
    expect(state.request).toBeNull()
    expect(state.pending).toBe('add')
  })

  it('keeps the pending operation when a unary applies to the second operand', () => {
    let state = press(initialState, ['1', '2', '.', '5', 'add', '9', 'sqrt'])
    expect(state.request?.continuation).toEqual({
      kind: 'unary-second',
      label: '√(9)'
    })
    state = succeed(state, 3)

    expect(state.pending).toBe('add')
    expect(state.expression).toBe('12.5 + √(9)')
    expect(state.current.text).toBe('3')

    state = press(state, ['equals'])
    expect(state.request?.operands).toEqual([12.5, 3])
    state = succeed(state, 15.5)
    expect(state.expression).toBe('12.5 + √(9) = 15.5')
  })

  it('ignores every math key while calculating', () => {
    const state = press(initialState, ['1', 'add', '2', 'equals'])
    for (const key of [
      '5',
      '.',
      'backspace',
      'sign',
      'add',
      'sqrt',
      'equals'
    ] as const) {
      expect(press(state, [key])).toBe(state)
    }
  })

  it('clears during a request but keeps the request sequence monotonic', () => {
    const state = press(initialState, ['1', 'add', '2', 'equals'])
    const cleared = press(state, ['clear'])

    expect(cleared.phase).toBe('idle')
    expect(cleared.current).toEqual({ kind: 'entry', text: '' })
    expect(cleared.sequence).toBe(state.sequence)

    const next = press(cleared, ['3', 'add', '4', 'equals'])
    expect(next.request?.id).toBeGreaterThan(state.request?.id ?? 0)
  })

  it('freezes the display and clears pending on a domain error', () => {
    let state = press(initialState, ['1', '0', 'divide', '0', 'equals'])
    state = fail(state, 'DIVISION_BY_ZERO')

    expect(state.phase).toBe('error')
    expect(state.error?.copy).toBe('Cannot divide by zero.')
    expect(state.pending).toBeNull()
    expect(state.current).toEqual({ kind: 'value', value: 0, text: '0' })
    expect(state.expression).toBe('10 ÷ 0')
    expect(press(state, ['equals'])).toBe(state)
    expect(press(state, ['backspace'])).toBe(state)
    expect(press(state, ['sign'])).toBe(state)

    const fresh = press(state, ['5'])
    expect(fresh.current).toEqual({ kind: 'entry', text: '5' })
    expect(fresh.error).toBeNull()

    const continued = press(state, ['add'])
    expect(continued.pending).toBe('add')
    expect(continued.expression).toBe('0 +')
  })

  it('retains the expression and retries the exact plan on a transient error', () => {
    let state = press(initialState, ['1', 'add', '2', 'equals'])
    const operation = state.request?.operation
    const operands = state.request?.operands
    state = fail(state, 'NETWORK_ERROR')

    expect(state.phase).toBe('error')
    expect(state.pending).toBe('add')
    expect(state.retry).not.toBeNull()
    expect(state.expression).toBe('1 + 2')

    state = press(state, ['equals'])
    expect(state.phase).toBe('calculating')
    expect(state.request?.operation).toBe(operation)
    expect(state.request?.operands).toEqual(operands)

    state = succeed(state, 3)
    expect(state.phase).toBe('result')
    expect(state.expression).toBe('1 + 2 = 3')
  })

  it('ignores stale responses', () => {
    const state = press(initialState, ['1', 'add', '2', 'equals'])
    const stale = reducer(state, {
      type: 'requestSucceeded',
      id: (state.request?.id ?? 0) - 1,
      result: 99
    })
    expect(stale).toBe(state)
  })

  it('carries the formatted result into the next request (WYSIWYG)', () => {
    let state = press(initialState, [
      '0',
      '.',
      '1',
      'add',
      '0',
      '.',
      '2',
      'equals'
    ])
    state = succeed(state, 0.30000000000000004)
    expect(state.current.text).toBe('0.3')

    state = press(state, ['multiply', '3', 'equals'])
    expect(state.request?.operation).toBe('multiply')
    expect(state.request?.operands).toEqual([0.3, 3])
  })

  it('treats power as a binary operation', () => {
    let state = press(initialState, ['2', 'power', '1', '0', 'equals'])
    expect(state.request?.operation).toBe('power')
    expect(state.request?.operands).toEqual([2, 10])

    state = succeed(state, 1024)
    expect(state.expression).toBe('2 ^ 10 = 1024')
  })

  it('starts a fresh calculation after a result', () => {
    let state = press(initialState, ['2', 'add', '3', 'equals'])
    state = succeed(state, 5)
    const fresh = press(state, ['7'])

    expect(fresh.phase).toBe('idle')
    expect(fresh.current).toEqual({ kind: 'entry', text: '7' })
    expect(fresh.first).toBeNull()
    expect(fresh.pending).toBeNull()
  })

  it('continues from a result with a new operator', () => {
    let state = press(initialState, ['2', 'add', '3', 'equals'])
    state = succeed(state, 5)
    const continued = press(state, ['subtract'])

    expect(continued.pending).toBe('subtract')
    expect(continued.first).toEqual({ value: 5, text: '5' })
    expect(continued.expression).toBe('5 −')
  })
})
