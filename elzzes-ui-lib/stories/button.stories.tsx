import type { Meta, StoryObj } from '@storybook/react-vite'
import { Button } from '../src'

const meta = {
  title: 'Components/Button',
  component: Button,
  args: {
    children: 'Button'
  },
  argTypes: {
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'outline', 'ghost', 'accent', 'danger']
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg', 'xl']
    }
  }
} satisfies Meta<typeof Button>

export default meta
type Story = StoryObj<typeof meta>

export const Primary: Story = {
  args: { variant: 'primary' }
}

export const AllVariants: Story = {
  render: (args) => (
    <div className="flex flex-wrap gap-2">
      {(
        [
          'primary',
          'secondary',
          'outline',
          'ghost',
          'accent',
          'danger'
        ] as const
      ).map((variant) => (
        <Button {...args} key={variant} variant={variant}>
          {variant}
        </Button>
      ))}
    </div>
  )
}

export const ExtraLarge: Story = {
  args: { size: 'xl' }
}
