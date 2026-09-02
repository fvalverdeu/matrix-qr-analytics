export type QrRequest = {
  matrix: number[][]
}

export type QrStatistics = {
  max: number
  min: number
  average: number
  sum: number
  hasDiagonalMatrix: boolean
}

export type QrResponse = {
  q: number[][]
  r: number[][]
  statistics: QrStatistics
}

export type ApiErrorEnvelope = {
  error: {
    code: string
    message: string
  }
}

export function isApiErrorEnvelope(value: unknown): value is ApiErrorEnvelope {
  if (!value || typeof value !== 'object') {
    return false
  }

  const candidate = value as { error?: unknown }
  if (!candidate.error || typeof candidate.error !== 'object') {
    return false
  }

  const error = candidate.error as { code?: unknown; message?: unknown }
  return typeof error.code === 'string' && typeof error.message === 'string'
}

function isNumericMatrix(value: unknown): value is number[][] {
  if (!Array.isArray(value) || value.length === 0) {
    return false
  }

  if (!value.every((row) => Array.isArray(row) && row.length > 0)) {
    return false
  }

  const firstRowLength = (value[0] as unknown[]).length
  if (firstRowLength === 0) {
    return false
  }

  return value.every((row) => {
    const typedRow = row as unknown[]
    if (typedRow.length !== firstRowLength) {
      return false
    }

    return typedRow.every(
      (cell) => typeof cell === 'number' && Number.isFinite(cell),
    )
  })
}

export function isQrResponse(value: unknown): value is QrResponse {
  if (!value || typeof value !== 'object') {
    return false
  }

  const candidate = value as {
    q?: unknown
    r?: unknown
    statistics?: unknown
  }

  if (!isNumericMatrix(candidate.q) || !isNumericMatrix(candidate.r)) {
    return false
  }

  if (!candidate.statistics || typeof candidate.statistics !== 'object') {
    return false
  }

  const stats = candidate.statistics as Record<string, unknown>
  return (
    typeof stats.max === 'number' &&
    Number.isFinite(stats.max) &&
    typeof stats.min === 'number' &&
    Number.isFinite(stats.min) &&
    typeof stats.average === 'number' &&
    Number.isFinite(stats.average) &&
    typeof stats.sum === 'number' &&
    Number.isFinite(stats.sum) &&
    typeof stats.hasDiagonalMatrix === 'boolean'
  )
}
