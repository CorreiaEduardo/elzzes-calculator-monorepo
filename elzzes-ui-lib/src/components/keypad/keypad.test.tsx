import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Button } from '../button'
import { Keypad } from '.'

describe('Keypad', () => {
  it('renders the accessory and main key grids', () => {
    render(
      <Keypad
        accessory={<Button>Square root</Button>}
        role="group"
        aria-label="Keys"
      >
        <Button>7</Button>
        <Button>8</Button>
      </Keypad>
    )

    const group = screen.getByRole('group', { name: 'Keys' })
    expect(group).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Square root' })
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '7' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '8' })).toBeInTheDocument()
  })

  it('renders without an accessory strip', () => {
    render(
      <Keypad role="group" aria-label="Digits">
        <Button>1</Button>
      </Keypad>
    )

    expect(screen.getByRole('group', { name: 'Digits' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '1' })).toBeInTheDocument()
  })
})
