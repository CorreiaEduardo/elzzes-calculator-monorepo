import type { HTMLAttributes } from 'react'

export type CardElement = 'div' | 'section' | 'article'
export type CardPadding = 'none' | 'sm' | 'md' | 'lg'

export interface CardProps extends HTMLAttributes<HTMLElement> {
  as?: CardElement
  padding?: CardPadding
}
