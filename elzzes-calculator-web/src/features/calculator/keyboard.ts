import type { CalculatorKey } from './types'

export function resolveKey(event: KeyboardEvent): CalculatorKey | null {
  if (
    event.defaultPrevented ||
    event.ctrlKey ||
    event.metaKey ||
    event.altKey ||
    event.isComposing
  ) {
    return null
  }

  const target = event.target as HTMLElement | null
  if (
    target &&
    (target.isContentEditable ||
      target.tagName === 'INPUT' ||
      target.tagName === 'TEXTAREA' ||
      target.tagName === 'SELECT')
  ) {
    return null
  }

  if (
    (event.key === 'Enter' || event.key === ' ') &&
    document.activeElement?.tagName === 'BUTTON'
  ) {
    return null
  }

  if (/^[0-9]$/.test(event.key)) {
    return event.key as CalculatorKey
  }

  switch (event.key) {
    case '.':
      return '.'
    case '+':
      return 'add'
    case '-':
      return 'subtract'
    case '*':
      return 'multiply'
    case '/':
      return 'divide'
    case '^':
      return 'power'
    case '%':
      return 'percent'
    case 'r':
    case 'R':
      return 'sqrt'
    case 's':
    case 'S':
      return 'sign'
    case 'Enter':
    case '=':
      return 'equals'
    case 'Backspace':
      return 'backspace'
    case 'Escape':
    case 'Delete':
      return 'clear'
    default:
      return null
  }
}
