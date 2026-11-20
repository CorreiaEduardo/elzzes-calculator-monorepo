import { z } from 'zod'
import { serverErrorCodes, type ServerErrorCode } from '../types'

export const operationSchema = z.enum([
  'add',
  'subtract',
  'multiply',
  'divide',
  'power',
  'sqrt',
  'percent'
])

const finiteNumber = z.number().refine(Number.isFinite, 'must be finite')

export const calculateResponseSchema = z.strictObject({
  operation: operationSchema,
  operands: z.array(finiteNumber),
  result: finiteNumber
})

export const errorCodeSchema = z.enum(serverErrorCodes)

export const calculateErrorSchema = z.strictObject({
  error: z.strictObject({
    code: errorCodeSchema,
    message: z.string()
  })
})

export type CalculateResponse = z.infer<typeof calculateResponseSchema>
export type CalculateErrorResponse = z.infer<typeof calculateErrorSchema>
export type { ServerErrorCode }
