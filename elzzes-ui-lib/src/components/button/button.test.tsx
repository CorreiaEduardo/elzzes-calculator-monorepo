import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { Button } from '.'

describe('Button', () => {
  it('renders a real button with its accessible name', () => {
    render(<Button>Save</Button>)

    const button = screen.getByRole('button', { name: 'Save' })
    expect(button).toHaveAttribute('type', 'button')
  })

  it('applies variant and size classes', () => {
    render(
      <Button variant="accent" size="xl">
        Advanced
      </Button>
    )

    const button = screen.getByRole('button', { name: 'Advanced' })
    expect(button.className).toContain('text-ui-accent')
    expect(button.className).toContain('min-h-14')
  })

  it('forwards aria attributes and handles activation', async () => {
    const onClick = vi.fn()
    render(
      <Button aria-keyshortcuts="+" onClick={onClick}>
        Add
      </Button>
    )

    const button = screen.getByRole('button', { name: 'Add' })
    expect(button).toHaveAttribute('aria-keyshortcuts', '+')
    await userEvent.click(button)
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('does not submit forms by default', () => {
    render(<Button>Action</Button>)
    expect(screen.getByRole('button', { name: 'Action' })).toHaveAttribute(
      'type',
      'button'
    )
  })
})
