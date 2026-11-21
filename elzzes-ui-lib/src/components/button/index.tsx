import { cva } from 'class-variance-authority'
import { cn } from '../../lib/cn'
import type { ButtonProps } from './button.types'

const button = cva(
  [
    'inline-flex items-center justify-center',
    'select-none touch-manipulation rounded-ui-key font-medium',
    'transition-colors duration-100 motion-reduce:transition-none',
    'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ui-accent',
    'disabled:pointer-events-none disabled:opacity-50'
  ],
  {
    variants: {
      variant: {
        primary: 'bg-ui-accent text-ui-accent-ink hover:bg-ui-accent-hover',
        secondary: 'bg-ui-key text-ui-ink hover:bg-ui-key-hover',
        outline: 'bg-ui-operator text-ui-ink hover:bg-ui-operator-hover',
        ghost: 'bg-transparent text-ui-ink hover:bg-ui-key',
        accent: 'bg-transparent text-ui-accent hover:bg-ui-accent-soft',
        danger: 'bg-transparent text-ui-danger hover:bg-ui-danger-soft'
      },
      size: {
        sm: 'px-3 py-1.5 text-sm',
        md: 'px-4 py-2 text-sm',
        lg: 'px-5 py-3 text-base',
        xl: 'min-h-14 px-4 text-xl'
      }
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md'
    }
  }
)

export function Button({
  variant = 'primary',
  size = 'md',
  type = 'button',
  className = '',
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      type={type}
      className={cn(button({ variant, size }), className)}
      {...props}
    >
      {children}
    </button>
  )
}
