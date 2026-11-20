import { displayEntry } from './entry'
import type {
  BinaryOperation,
  CalculatorState,
  Operand,
  UnaryOperation
} from './types'

export const binarySymbols: Record<BinaryOperation, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: '^'
}

export function unaryLabel(operation: UnaryOperation, source: string): string {
  return operation === 'sqrt' ? `√(${source})` : `${source}%`
}

export function operandText(operand: Operand): string {
  return operand.kind === 'entry' ? displayEntry(operand.text) : operand.text
}

export function operandValue(operand: Operand): number {
  return operand.kind === 'entry'
    ? Number(displayEntry(operand.text))
    : operand.value
}

export function displayValue(state: CalculatorState): string {
  return operandText(state.current)
}

export function expressionValue(state: CalculatorState, busy: string): string {
  return state.phase === 'calculating' ? busy : state.expression
}
