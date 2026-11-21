const SIGNIFICANT_DIGITS = 12
const EXPONENTIAL_HIGH = 1e12
const EXPONENTIAL_LOW = 1e-9

export function formatNumber(value: number): string {
  if (!Number.isFinite(value)) {
    return '0'
  }
  if (value === 0) {
    return '0'
  }

  const rounded = Number(value.toPrecision(SIGNIFICANT_DIGITS))
  const absolute = Math.abs(rounded)

  if (absolute >= EXPONENTIAL_HIGH || absolute < EXPONENTIAL_LOW) {
    return formatExponential(rounded)
  }

  const exponent = Math.floor(Math.log10(absolute))
  const decimals = Math.max(0, SIGNIFICANT_DIGITS - 1 - exponent)
  return trimDecimal(rounded.toFixed(decimals))
}

function formatExponential(value: number): string {
  const [mantissa, exponent] = value
    .toExponential(SIGNIFICANT_DIGITS - 1)
    .split('e')
  const trimmed = mantissa.replace(/\.?0+$/, '')
  const exponentValue = Number(exponent)
  const sign = exponentValue >= 0 ? '+' : '-'
  return `${trimmed}e${sign}${Math.abs(exponentValue)}`
}

function trimDecimal(value: string): string {
  return value.includes('.') ? value.replace(/\.?0+$/, '') : value
}
