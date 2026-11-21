import { useCallback, useEffect, useReducer, useRef } from 'react'
import { useMutation } from '@tanstack/react-query'
import { calculate } from '../api/calculate'
import { errorForCode } from '../copy'
import { reducer } from '../reducer'
import { initialState, type ActiveRequest, type CalculatorKey } from '../types'

export function useCalculator() {
  const [state, dispatch] = useReducer(reducer, initialState)
  const startedRequest = useRef<number | null>(null)
  const abortRef = useRef<AbortController | null>(null)

  const mutation = useMutation({
    mutationFn: async ({
      request,
      signal
    }: {
      request: ActiveRequest
      signal: AbortSignal
    }) => ({
      id: request.id,
      outcome: await calculate(request.operation, request.operands, signal)
    }),
    retry: false,
    networkMode: 'always',
    onSuccess: ({ id, outcome }) => {
      if (outcome.kind === 'success') {
        dispatch({ type: 'requestSucceeded', id, result: outcome.result })
      } else if (outcome.kind === 'error') {
        dispatch({ type: 'requestFailed', id, error: outcome.error })
      }
    },
    onError: (_error, variables) => {
      dispatch({
        type: 'requestFailed',
        id: variables.request.id,
        error: errorForCode('INTERNAL_ERROR')
      })
    }
  })

  const { mutate } = mutation

  useEffect(() => {
    const request = state.request
    if (!request || startedRequest.current === request.id) {
      return
    }
    startedRequest.current = request.id
    const controller = new AbortController()
    abortRef.current = controller
    mutate({ request, signal: controller.signal })
  }, [state.request, mutate])

  const press = useCallback((key: CalculatorKey) => {
    if (key === 'clear') {
      abortRef.current?.abort()
      abortRef.current = null
    }
    dispatch({ type: 'press', key })
  }, [])

  return { state, press }
}
