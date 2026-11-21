export const MAX_ENTRY_DIGITS = 15

export function digitCount(text: string): number {
  return text.replace(/[-.]/g, '').length
}

export function displayEntry(text: string): string {
  return text === '' ? '0' : text
}

export function entryValue(text: string): number {
  return text === '' || text === '-' ? 0 : Number(text)
}

export function inputDigit(text: string, digit: string): string {
  if (text === '' || text === '0') {
    return digit === '0' ? '0' : digit
  }
  if (digitCount(text) >= MAX_ENTRY_DIGITS) {
    return text
  }
  return text + digit
}

export function inputDot(text: string): string {
  if (text.includes('.')) {
    return text
  }
  if (text === '' || text === '0') {
    return '0.'
  }
  return text + '.'
}

export function backspace(text: string): string {
  const next = text.slice(0, -1)
  if (next === '' || next === '-') {
    return ''
  }
  return next
}

export function toggleSign(text: string): string {
  if (text === '' || Number(text) === 0) {
    return text
  }
  return text.startsWith('-') ? text.slice(1) : `-${text}`
}
