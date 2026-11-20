import { describe, expect, it } from 'vitest'
import { formatNumber } from './format'

describe('formatNumber', () => {
  it('trims trailing zeros and keeps up to 12 significant digits', () => {
    expect(formatNumber(20)).toBe('20')
    expect(formatNumber(0.125)).toBe('0.125')
    expect(formatNumber(0.1 + 0.2)).toBe('0.3')
    expect(formatNumber(1 / 3)).toBe('0.333333333333')
    expect(formatNumber(-12.5)).toBe('-12.5')
  })

  it('uses exponential notation at the defined thresholds', () => {
    expect(formatNumber(1e12)).toBe('1e+12')
    expect(formatNumber(999999999999)).toBe('999999999999')
    expect(formatNumber(999999999999.9)).toBe('1e+12')
    expect(formatNumber(1e-9)).toBe('0.000000001')
    expect(formatNumber(1e-10)).toBe('1e-10')
    expect(formatNumber(1.2e-10)).toBe('1.2e-10')
    expect(formatNumber(1234567890123)).toBe('1.23456789012e+12')
  })

  it('renders zero and negative zero as 0', () => {
    expect(formatNumber(0)).toBe('0')
    expect(formatNumber(-0)).toBe('0')
  })

  it('formats large finite values without separators', () => {
    expect(formatNumber(Number.MAX_VALUE)).toMatch(/^1\.79769313486e\+308$/)
  })
})
