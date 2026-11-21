import type { ButtonVariant } from '@elzzes/ui'
import type { CalculatorKey } from './types'

export interface KeyDefinition {
  key: CalculatorKey
  label: string
  name?: string
  variant: ButtonVariant
  shortcut?: string
  span?: 2
}

export const advancedKeys: KeyDefinition[] = [
  {
    key: 'sqrt',
    label: '√',
    name: 'square root',
    variant: 'accent',
    shortcut: 'r'
  },
  {
    key: 'power',
    label: 'xʸ',
    name: 'power',
    variant: 'accent',
    shortcut: '^'
  },
  {
    key: 'percent',
    label: '%',
    name: 'percent',
    variant: 'accent',
    shortcut: '%'
  }
]

export const mainKeys: KeyDefinition[] = [
  {
    key: 'clear',
    label: 'AC',
    name: 'clear',
    variant: 'danger',
    shortcut: 'Escape Delete'
  },
  {
    key: 'backspace',
    label: '⌫',
    name: 'backspace',
    variant: 'ghost',
    shortcut: 'Backspace'
  },
  {
    key: 'sign',
    label: '±',
    name: 'toggle sign',
    variant: 'ghost',
    shortcut: 's'
  },
  {
    key: 'divide',
    label: '÷',
    name: 'divide',
    variant: 'outline',
    shortcut: '/'
  },

  { key: '7', label: '7', variant: 'secondary', shortcut: '7' },
  { key: '8', label: '8', variant: 'secondary', shortcut: '8' },
  { key: '9', label: '9', variant: 'secondary', shortcut: '9' },
  {
    key: 'multiply',
    label: '×',
    name: 'multiply',
    variant: 'outline',
    shortcut: '*'
  },

  { key: '4', label: '4', variant: 'secondary', shortcut: '4' },
  { key: '5', label: '5', variant: 'secondary', shortcut: '5' },
  { key: '6', label: '6', variant: 'secondary', shortcut: '6' },
  {
    key: 'subtract',
    label: '−',
    name: 'subtract',
    variant: 'outline',
    shortcut: '-'
  },

  { key: '1', label: '1', variant: 'secondary', shortcut: '1' },
  { key: '2', label: '2', variant: 'secondary', shortcut: '2' },
  { key: '3', label: '3', variant: 'secondary', shortcut: '3' },
  { key: 'add', label: '+', name: 'add', variant: 'outline', shortcut: '+' },

  { key: '0', label: '0', variant: 'secondary', shortcut: '0', span: 2 },
  {
    key: '.',
    label: '.',
    name: 'decimal point',
    variant: 'secondary',
    shortcut: '.'
  },
  {
    key: 'equals',
    label: '=',
    name: 'equals',
    variant: 'primary',
    shortcut: 'Enter'
  }
]
