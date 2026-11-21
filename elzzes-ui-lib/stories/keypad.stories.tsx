import type { Meta, StoryObj } from '@storybook/react-vite'
import { Button, Keypad } from '../src'

const meta = {
  title: 'Components/Keypad',
  component: Keypad,
  args: {
    children: null
  }
} satisfies Meta<typeof Keypad>

export default meta
type Story = StoryObj<typeof meta>

const keys = ['7', '8', '9', '÷', '4', '5', '6', '×', '1', '2', '3', '−']

export const Default: Story = {
  render: () => (
    <Keypad
      role="group"
      aria-label="Example keypad"
      accessory={
        <>
          <Button variant="accent" size="xl">
            √
          </Button>
          <Button variant="accent" size="xl">
            xʸ
          </Button>
          <Button variant="accent" size="xl">
            %
          </Button>
        </>
      }
    >
      {keys.map((key) => (
        <Button key={key} variant="secondary" size="xl">
          {key}
        </Button>
      ))}
    </Keypad>
  )
}
