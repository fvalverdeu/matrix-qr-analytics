import { describe, expect, it } from 'vitest'
import { formatNumberForDisplay } from './numberFormat'

describe('formatNumberForDisplay', () => {
  it('keeps integers clean', () => {
    expect(formatNumberForDisplay(1)).toBe('1')
  })

  it('keeps meaningful decimals', () => {
    expect(formatNumberForDisplay(1.5)).toBe('1.5')
  })

  it('rounds values with more than 6 decimals', () => {
    expect(formatNumberForDisplay(0.3333333333)).toBe('0.333333')
  })

  it('removes trailing zeros after rounding', () => {
    expect(formatNumberForDisplay(1.23)).toBe('1.23')
  })

  it('normalizes rounded negative zero to 0', () => {
    expect(formatNumberForDisplay(-0.00000000001)).toBe('0')
  })

  it('displays zero as 0', () => {
    expect(formatNumberForDisplay(0)).toBe('0')
  })
})
