import { Button } from '@elzzes/ui'
import type { KeyDefinition } from '../keys'
import type { CalculatorKey } from '../types'

interface KeyButtonProps {
  definition: KeyDefinition
  onPress: (key: CalculatorKey) => void
}

export function KeyButton({ definition, onPress }: KeyButtonProps) {
  const { key, label, name, variant, shortcut, span } = definition

  return (
    <Button
      variant={variant}
      size="xl"
      className={span === 2 ? 'col-span-2 w-full' : 'w-full'}
      aria-label={name}
      aria-keyshortcuts={shortcut}
      onClick={() => onPress(key)}
    >
      {label}
    </Button>
  )
}
