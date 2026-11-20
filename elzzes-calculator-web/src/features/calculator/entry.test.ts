import { describe, expect, it } from 'vitest'
import {
  MAX_ENTRY_DIGITS,
  backspace,
  displayEntry,
  entryValue,
  inputDigit,
  inputDot,
  toggleSign
} from './entry'

describe('entry rules', () => {
  it('replaces leading zeros and appends digits', () => {
    expect(inputDigit('', '5')).toBe('5')
    expect(inputDigit('', '0')).toBe('0')
    expect(inputDigit('0', '5')).toBe('5')
    expect(inputDigit('0', '0')).toBe('0')
    expect(inputDigit('1', '2')).toBe('12')
    expect(inputDigit('0.', '5')).toBe('0.5')
  })

  it('enforces the 15-digit cap excluding sign and dot', () => {
    const full = '123456789012345'
    expect(full).toHaveLength(MAX_ENTRY_DIGITS)
    expect(inputDigit(full, '6')).toBe(full)
    expect(inputDigit('-1.23456789012345', '6')).toBe('-1.23456789012345')
  })

  it('allows one decimal point', () => {
    expect(inputDot('')).toBe('0.')
    expect(inputDot('0')).toBe('0.')
    expect(inputDot('1')).toBe('1.')
    expect(inputDot('1.5')).toBe('1.5')
  })

  it('backspaces entries and normalizes empty values', () => {
    expect(backspace('5')).toBe('')
    expect(backspace('-5')).toBe('')
    expect(backspace('0.')).toBe('0')
    expect(backspace('-')).toBe('')
    expect(backspace('')).toBe('')
  })

  it('toggles sign only for non-zero entries', () => {
    expect(toggleSign('5')).toBe('-5')
    expect(toggleSign('-5')).toBe('5')
    expect(toggleSign('')).toBe('')
    expect(toggleSign('0')).toBe('0')
    expect(toggleSign('0.')).toBe('0.')
  })

  it('projects empty entries and values', () => {
    expect(displayEntry('')).toBe('0')
    expect(displayEntry('12')).toBe('12')
    expect(entryValue('')).toBe(0)
    expect(entryValue('12.5')).toBe(12.5)
  })
})
