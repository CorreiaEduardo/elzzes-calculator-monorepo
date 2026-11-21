import type { Meta, StoryObj } from '@storybook/react-vite'
import { Display } from '../src'

const meta = {
  title: 'Components/Display',
  component: Display,
  args: {
    value: '12.5',
    expression: '12.5 +'
  }
} satisfies Meta<typeof Display>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Busy: Story = {
  args: { busy: true, expression: 'Calculating…' }
}

export const WithError: Story = {
  args: { value: '0', expression: '10 ÷ 0', error: 'Cannot divide by zero.' }
}

export const LongValue: Story = {
  args: { value: '0.333333333333', expression: '1 ÷ 3' }
}
