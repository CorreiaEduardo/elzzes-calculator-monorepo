import { cn } from '../../lib/cn'
import type { CardPadding, CardProps } from './card.types'

const paddingClasses: Record<CardPadding, string> = {
  none: '',
  sm: 'p-3',
  md: 'p-4',
  lg: 'p-6'
}

export function Card({
  as: Component = 'div',
  padding = 'md',
  className = '',
  children,
  ...props
}: CardProps) {
  return (
    <Component
      className={cn(
        'rounded-ui-card border border-ui-border bg-ui-surface text-ui-ink shadow-ui-card',
        paddingClasses[padding],
        className
      )}
      {...props}
    >
      {children}
    </Component>
  )
}
