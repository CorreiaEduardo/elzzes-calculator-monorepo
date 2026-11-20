import type { Meta, StoryObj } from '@storybook/react-vite'
import { Card } from '../src'

const meta = {
  title: 'Components/Card',
  component: Card,
  args: {
    children: 'Card content',
    padding: 'md'
  }
} satisfies Meta<typeof Card>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Section: Story = {
  args: { as: 'section', padding: 'lg', 'aria-label': 'Example section' }
}
