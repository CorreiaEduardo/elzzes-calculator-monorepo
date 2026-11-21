import type { HTMLAttributes } from 'react'

export interface DisplayProps extends HTMLAttributes<HTMLDivElement> {
  value: string
  expression?: string
  error?: string | null
  busy?: boolean
}
