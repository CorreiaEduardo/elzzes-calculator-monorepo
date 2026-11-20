import { afterEach, describe, expect, it } from 'vitest'
import { resolveKey } from './keyboard'

function keydown(key: string, init: KeyboardEventInit = {}): KeyboardEvent {
  return new KeyboardEvent('keydown', { key, ...init })
}

afterEach(() => {
  document.body.innerHTML = ''
})

describe('resolveKey', () => {
  it('maps digits, operators, unary keys, and control keys', () => {
    expect(resolveKey(keydown('7'))).toBe('7')
    expect(resolveKey(keydown('.'))).toBe('.')
    expect(resolveKey(keydown('+'))).toBe('add')
    expect(resolveKey(keydown('-'))).toBe('subtract')
    expect(resolveKey(keydown('*'))).toBe('multiply')
    expect(resolveKey(keydown('/'))).toBe('divide')
    expect(resolveKey(keydown('^'))).toBe('power')
    expect(resolveKey(keydown('%'))).toBe('percent')
    expect(resolveKey(keydown('r'))).toBe('sqrt')
    expect(resolveKey(keydown('s'))).toBe('sign')
    expect(resolveKey(keydown('Enter'))).toBe('equals')
    expect(resolveKey(keydown('Backspace'))).toBe('backspace')
    expect(resolveKey(keydown('Escape'))).toBe('clear')
    expect(resolveKey(keydown('Delete'))).toBe('clear')
  })

  it('ignores modifiers, composed events, and already-handled events', () => {
    expect(resolveKey(keydown('1', { ctrlKey: true }))).toBeNull()
    expect(resolveKey(keydown('1', { metaKey: true }))).toBeNull()
    expect(resolveKey(keydown('1', { altKey: true }))).toBeNull()
    expect(resolveKey(keydown('1', { isComposing: true }))).toBeNull()
    expect(resolveKey(keydown('x'))).toBeNull()
  })

  it('ignores text inputs', () => {
    const input = document.createElement('input')
    const event = keydown('1')
    input.dispatchEvent(event)
    expect(resolveKey(event)).toBeNull()
  })

  it('lets a focused button handle Enter natively', () => {
    const button = document.createElement('button')
    document.body.append(button)
    button.focus()
    const event = keydown('Enter')
    button.dispatchEvent(event)
    expect(resolveKey(event)).toBeNull()
  })
})
