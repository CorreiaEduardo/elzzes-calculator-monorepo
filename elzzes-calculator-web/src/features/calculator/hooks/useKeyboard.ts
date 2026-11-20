import { useEffect } from 'react'
import { resolveKey } from '../keyboard'
import type { CalculatorKey } from '../types'

export function useKeyboard(press: (key: CalculatorKey) => void) {
  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      const key = resolveKey(event)
      if (!key) {
        return
      }
      event.preventDefault()
      press(key)
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [press])
}
