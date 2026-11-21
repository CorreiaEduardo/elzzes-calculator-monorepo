import { http, HttpResponse } from 'msw'
import { describe, expect, it, vi } from 'vitest'
import { server } from '../../../test/server'
import { calculate } from './calculate'

const endpoint = '*/api/v1/calculate'

function successResponse(result: number) {
  return HttpResponse.json({
    operation: 'add',
    operands: [1, 2],
    result
  })
}

describe('calculate API client', () => {
  it('posts the operation and operands and accepts a valid response', async () => {
    let received: unknown
    server.use(
      http.post(endpoint, async ({ request }) => {
        received = await request.json()
        return successResponse(3)
      })
    )

    const outcome = await calculate('add', [1, 2])

    expect(received).toEqual({ operation: 'add', operands: [1, 2] })
    expect(outcome).toEqual({ kind: 'success', result: 3 })
  })

  it('maps 422 domain codes to their copy', async () => {
    server.use(
      http.post(endpoint, () =>
        HttpResponse.json(
          {
            error: {
              code: 'DIVISION_BY_ZERO',
              message: 'cannot divide by zero'
            }
          },
          { status: 422 }
        )
      )
    )

    const outcome = await calculate('divide', [10, 0])

    expect(outcome).toMatchObject({
      kind: 'error',
      error: {
        kind: 'domain',
        code: 'DIVISION_BY_ZERO',
        copy: 'Cannot divide by zero.'
      }
    })
  })

  it('maps documented 4xx validation errors to the invalid-calculation copy', async () => {
    server.use(
      http.post(endpoint, () =>
        HttpResponse.json(
          { error: { code: 'INVALID_OPERAND', message: 'bad operand' } },
          { status: 400 }
        )
      )
    )

    const outcome = await calculate('add', [1, 2])

    expect(outcome).toMatchObject({
      kind: 'error',
      error: {
        kind: 'domain',
        code: 'INVALID_CALCULATION',
        copy: "That calculation isn't valid."
      }
    })
  })

  it('maps 5xx errors to a transient internal error', async () => {
    server.use(
      http.post(endpoint, () =>
        HttpResponse.json(
          { error: { code: 'INTERNAL_ERROR', message: 'internal error' } },
          { status: 500 }
        )
      )
    )

    const outcome = await calculate('add', [1, 2])

    expect(outcome).toMatchObject({
      kind: 'error',
      error: {
        kind: 'transient',
        code: 'INTERNAL_ERROR',
        copy: 'Something went wrong. Press = to retry.'
      }
    })
  })

  it('maps network failures to NETWORK_ERROR', async () => {
    server.use(http.post(endpoint, () => HttpResponse.error()))

    const outcome = await calculate('add', [1, 2])

    expect(outcome).toMatchObject({
      kind: 'error',
      error: { kind: 'transient', code: 'NETWORK_ERROR' }
    })
  })

  it('maps malformed success bodies and schema drift to INVALID_RESPONSE', async () => {
    server.use(http.post(endpoint, () => HttpResponse.text('not json')))
    expect(await calculate('add', [1, 2])).toMatchObject({
      kind: 'error',
      error: { code: 'INVALID_RESPONSE' }
    })

    server.use(
      http.post(endpoint, () =>
        HttpResponse.json({
          operation: 'add',
          operands: [1, 2],
          result: 3,
          extra: true
        })
      )
    )
    expect(await calculate('add', [1, 2])).toMatchObject({
      kind: 'error',
      error: { code: 'INVALID_RESPONSE' }
    })
  })

  it('maps malformed error envelopes to INVALID_RESPONSE', async () => {
    server.use(
      http.post(endpoint, () =>
        HttpResponse.json({ error: { code: 'NOT_A_CODE' } }, { status: 422 })
      )
    )

    expect(await calculate('add', [1, 2])).toMatchObject({
      kind: 'error',
      error: { code: 'INVALID_RESPONSE' }
    })
  })

  it('returns aborted for a pre-aborted signal', async () => {
    const controller = new AbortController()
    controller.abort()

    expect(await calculate('add', [1, 2], controller.signal)).toEqual({
      kind: 'aborted'
    })
  })

  it('times out after ten seconds', async () => {
    vi.useFakeTimers()
    server.use(http.post(endpoint, () => new Promise(() => {})))

    const pending = calculate('add', [1, 2])
    await vi.advanceTimersByTimeAsync(10_000)
    const outcome = await pending

    expect(outcome).toMatchObject({
      kind: 'error',
      error: { kind: 'transient', code: 'TIMEOUT' }
    })
  })
})
