import type { ButtonHTMLAttributes } from 'react'

export type ButtonVariant =
  'primary' | 'secondary' | 'outline' | 'ghost' | 'accent' | 'danger'

export type ButtonSize = 'sm' | 'md' | 'lg' | 'xl'

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
}
