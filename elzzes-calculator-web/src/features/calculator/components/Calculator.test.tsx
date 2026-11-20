import { StrictMode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import { server } from '../../../test/server'
import { Calculator } from './Calculator'

const endpoint = '*/api/v1/calculate'

function renderCalculator(strict = false) {
  const client = new QueryClient({
    defaultOptions: {
      mutations: { retry: false, networkMode: 'always' }
    }
  })
  const tree = (
    <QueryClientProvider client={client}>
      <Calculator />
    </QueryClientProvider>
  )
  return render(strict ? <StrictMode>{tree}</StrictMode> : tree)
}

function displayValue(container: HTMLElement): string {
  return container.querySelector('.tabular-nums')?.textContent ?? ''
}

function expressionText(container: HTMLElement): string {
  return container.querySelector('[aria-live="polite"]')?.textContent ?? ''
}

async function clickKeys(
  user: ReturnType<typeof userEvent.setup>,
  keys: string[]
) {
  for (const name of keys) {
    await user.click(screen.getByRole('button', { name }))
  }
}

describe('Calculator', () => {
  it('renders the keypad, display, and accessible names', () => {
    const { container } = renderCalculator()

    expect(
      screen.getByRole('region', { name: 'Calculator' })
    ).toBeInTheDocument()
    expect(
      screen.getByRole('group', { name: 'Calculator keypad' })
    ).toBeInTheDocument()

    for (const name of [
      'add',
      'subtract',
      'multiply',
      'divide',
      'power',
      'square root',
      'percent',
      'clear',
      'backspace',
      'toggle sign',
      'decimal point',
      'equals'
    ]) {
      expect(screen.getByRole('button', { name })).toBeInTheDocument()
    }

    expect(displayValue(container)).toBe('0')
    expect(expressionText(container)).toBe('')
  })

  it('trusts the server result and keeps focus on the pressed key', async () => {
    const user = userEvent.setup()
    let received: unknown

    server.use(
      http.post(endpoint, async ({ request }) => {
        received = await request.json()
        return HttpResponse.json({
          operation: 'add',
          operands: [12, 7],
          result: 999
        })
      })
    )

    const { container } = renderCalculator()
    await clickKeys(user, ['1', '2', 'add', '7'])

    const equals = screen.getByRole('button', { name: 'equals' })
    await user.click(equals)

    await waitFor(() => expect(displayValue(container)).toBe('999'))
    expect(expressionText(container)).toBe('12 + 7 = 999')
    expect(received).toEqual({ operation: 'add', operands: [12, 7] })
    expect(document.activeElement).toBe(equals)
  })

  it('shows domain errors and recovers without a retry', async () => {
    const user = userEvent.setup()
    let requests = 0

    server.use(
      http.post(endpoint, () => {
        requests += 1
        return HttpResponse.json(
          {
            error: {
              code: 'DIVISION_BY_ZERO',
              message: 'cannot divide by zero'
            }
          },
          { status: 422 }
        )
      })
    )

    const { container } = renderCalculator()
    await clickKeys(user, ['1', '0', 'divide', '0', 'equals'])

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Cannot divide by zero.'
    )
    expect(expressionText(container)).toBe('10 ÷ 0')
    expect(displayValue(container)).toBe('0')

    await user.click(screen.getByRole('button', { name: 'equals' }))
    expect(requests).toBe(1)
  })

  it('locks math keys while in flight and AC aborts', async () => {
    const user = userEvent.setup()
    let releaseRequest = () => {}

    server.use(
      http.post(endpoint, async () => {
        await new Promise<void>((resolve) => {
          releaseRequest = resolve
        })
        return HttpResponse.json({
          operation: 'add',
          operands: [1, 2],
          result: 3
        })
      })
    )

    const { container } = renderCalculator()
    await clickKeys(user, ['1', 'add', '2', 'equals'])

    await waitFor(() =>
      expect(
        screen.getByRole('region', { name: 'Calculator' })
      ).toHaveAttribute('aria-busy', 'true')
    )

    await user.click(screen.getByRole('button', { name: '3' }))
    expect(displayValue(container)).toBe('2')

    await user.click(screen.getByRole('button', { name: 'clear' }))
    expect(displayValue(container)).toBe('0')
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()

    releaseRequest()
    await waitFor(() => expect(displayValue(container)).toBe('0'))
  })

  it('retries a transient failure with the exact same request', async () => {
    const user = userEvent.setup()
    const bodies: unknown[] = []

    server.use(
      http.post(endpoint, async ({ request }) => {
        bodies.push(await request.json())
        if (bodies.length === 1) {
          return HttpResponse.json(
            { error: { code: 'INTERNAL_ERROR', message: 'internal error' } },
            { status: 500 }
          )
        }
        return HttpResponse.json({
          operation: 'add',
          operands: [1, 2],
          result: 3
        })
      })
    )

    const { container } = renderCalculator()
    await clickKeys(user, ['1', 'add', '2', 'equals'])

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Something went wrong. Press = to retry.'
    )

    await user.click(screen.getByRole('button', { name: 'equals' }))
    await waitFor(() => expect(displayValue(container)).toBe('3'))

    expect(bodies).toHaveLength(2)
    expect(bodies[1]).toEqual(bodies[0])
  })

  it('carries the formatted display value into the next operation', async () => {
    const user = userEvent.setup()
    let nextOperands: unknown

    server.use(
      http.post(endpoint, async ({ request }) => {
        const body = (await request.json()) as {
          operation: string
          operands: number[]
        }
        if (body.operation === 'multiply') {
          nextOperands = body.operands
          return HttpResponse.json({
            operation: 'multiply',
            operands: body.operands,
            result: 0.9
          })
        }
        return HttpResponse.json({
          operation: 'add',
          operands: body.operands,
          result: 0.30000000000000004
        })
      })
    )

    const { container } = renderCalculator()
    await clickKeys(user, [
      '0',
      'decimal point',
      '1',
      'add',
      '0',
      'decimal point',
      '2',
      'equals'
    ])
    await waitFor(() => expect(displayValue(container)).toBe('0.3'))

    await clickKeys(user, ['multiply', '3', 'equals'])
    await waitFor(() => expect(displayValue(container)).toBe('0.9'))
    expect(nextOperands).toEqual([0.3, 3])
  })

  it('supports keyboard input', async () => {
    const user = userEvent.setup()

    server.use(
      http.post(endpoint, async ({ request }) => {
        const body = (await request.json()) as {
          operation: string
          operands: number[]
        }
        return HttpResponse.json({
          operation: body.operation,
          operands: body.operands,
          result: 56
        })
      })
    )

    const { container } = renderCalculator()
    await user.keyboard('7*8{Enter}')

    await waitFor(() => expect(displayValue(container)).toBe('56'))
    expect(expressionText(container)).toBe('7 × 8 = 56')
  })

  it('issues one request under StrictMode', async () => {
    const user = userEvent.setup()
    let requests = 0

    server.use(
      http.post(endpoint, () => {
        requests += 1
        return HttpResponse.json({
          operation: 'add',
          operands: [1, 2],
          result: 3
        })
      })
    )

    const { container } = renderCalculator(true)
    await clickKeys(user, ['1', 'add', '2', 'equals'])
    await waitFor(() => expect(displayValue(container)).toBe('3'))

    expect(requests).toBe(1)
  })
})
