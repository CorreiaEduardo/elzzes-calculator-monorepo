import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Display } from '.'

describe('Display', () => {
  it('renders the value and a polite expression region', () => {
    render(<Display value="12.5" expression="12.5 +" />)

    expect(screen.getByText('12.5')).toBeInTheDocument()
    expect(screen.getByText('12.5 +')).toHaveAttribute('aria-live', 'polite')
  })

  it('renders an alert and busy state', () => {
    render(<Display value="0" error="Cannot divide by zero." busy />)

    expect(screen.getByRole('alert')).toHaveTextContent(
      'Cannot divide by zero.'
    )
    expect(screen.getByText('0').parentElement).toHaveAttribute(
      'aria-busy',
      'true'
    )
  })

  it('does not render an alert without an error', () => {
    render(<Display value="0" />)
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('steps the value type down as the value grows', () => {
    const { rerender } = render(<Display value="1234567890" />)
    expect(screen.getByText('1234567890').className).toContain('text-4xl')

    rerender(<Display value="1234567890123" />)
    expect(screen.getByText('1234567890123').className).toContain('text-3xl')

    rerender(<Display value="12345678901234567" />)
    expect(screen.getByText('12345678901234567').className).toContain(
      'text-2xl'
    )
  })
})
