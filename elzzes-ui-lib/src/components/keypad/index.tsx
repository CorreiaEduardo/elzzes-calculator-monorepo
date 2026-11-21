import { cn } from '../../lib/cn'
import type { KeypadProps } from './keypad.types'

export function Keypad({
  accessory,
  className = '',
  children,
  ...props
}: KeypadProps) {
  return (
    <div className={cn('flex flex-col gap-2', className)} {...props}>
      {accessory ? (
        <div className="grid grid-cols-3 gap-2">{accessory}</div>
      ) : null}
      <div className="grid grid-cols-4 gap-2">{children}</div>
    </div>
  )
}
