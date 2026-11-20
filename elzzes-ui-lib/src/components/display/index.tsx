import { cn } from '../../lib/cn'
import type { DisplayProps } from './display.types'

function valueSize(value: string): string {
  if (value.length > 16) {
    return 'text-2xl'
  }
  if (value.length > 12) {
    return 'text-3xl'
  }
  if (value.length > 9) {
    return 'text-4xl'
  }
  return 'text-5xl'
}

export function Display({
  value,
  expression = '',
  error = null,
  busy = false,
  className = '',
  ...props
}: DisplayProps) {
  return (
    <div
      className={cn('flex flex-col gap-2', className)}
      aria-busy={busy || undefined}
      {...props}
    >
      <p
        aria-live="polite"
        aria-atomic="true"
        className="min-h-6 overflow-hidden text-ellipsis whitespace-nowrap text-sm text-ui-muted"
      >
        {expression}
      </p>
      <div
        className={cn(
          'overflow-hidden whitespace-nowrap text-right font-semibold tabular-nums text-ui-ink',
          valueSize(value)
        )}
      >
        {value}
      </div>
      {error ? (
        <p role="alert" className="text-sm text-ui-danger">
          {error}
        </p>
      ) : null}
    </div>
  )
}
