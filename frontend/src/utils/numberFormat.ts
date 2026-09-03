export function formatNumberForDisplay(value: number): string {
  if (!Number.isFinite(value)) {
    return '0'
  }

  const rounded = Number(value.toFixed(6))

  if (Object.is(rounded, -0) || rounded === 0) {
    return '0'
  }

  return rounded.toString()
}
