import { Card, Display, Keypad } from '@elzzes/ui'
import { busyText } from '../copy'
import { displayValue, expressionValue } from '../expression'
import { useCalculator } from '../hooks/useCalculator'
import { useKeyboard } from '../hooks/useKeyboard'
import { advancedKeys, mainKeys } from '../keys'
import { KeyButton } from './KeyButton'

export function Calculator() {
  const { state, press } = useCalculator()
  useKeyboard(press)

  const busy = state.phase === 'calculating'

  return (
    <Card
      as="section"
      aria-label="Calculator"
      aria-busy={busy || undefined}
      padding="md"
      className="w-full max-w-[420px]"
    >
      <Display
        value={displayValue(state)}
        expression={expressionValue(state, busyText)}
        error={state.error?.copy ?? null}
        busy={busy}
      />
      <Keypad
        className="mt-4"
        role="group"
        aria-label="Calculator keypad"
        accessory={advancedKeys.map((definition) => (
          <KeyButton
            key={definition.key}
            definition={definition}
            onPress={press}
          />
        ))}
      >
        {mainKeys.map((definition) => (
          <KeyButton
            key={definition.key}
            definition={definition}
            onPress={press}
          />
        ))}
      </Keypad>
    </Card>
  )
}
