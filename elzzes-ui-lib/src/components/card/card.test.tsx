import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Card } from '.'

describe('Card', () => {
  it('renders a surface with its content', () => {
    render(<Card aria-label="Surface">Content</Card>)

    expect(screen.getByLabelText('Surface')).toHaveTextContent('Content')
  })

  it('supports semantic elements and padding', () => {
    render(
      <Card as="section" padding="lg" aria-label="Panel">
        Body
      </Card>
    )

    const panel = screen.getByRole('region', { name: 'Panel' })
    expect(panel.tagName).toBe('SECTION')
    expect(panel.className).toContain('bg-ui-surface')
    expect(panel.className).toContain('p-6')
  })
})
