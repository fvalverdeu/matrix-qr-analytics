import { apiBaseUrl } from '../config/env'
import {
  isApiErrorEnvelope,
  isQrResponse,
  type QrRequest,
  type QrResponse,
} from '../types/qr'

export type QrApiErrorKind = 'network' | 'http' | 'invalid-response'

export class QrApiError extends Error {
  readonly kind: QrApiErrorKind
  readonly status?: number
  readonly code?: string

  constructor(
    kind: QrApiErrorKind,
    message: string,
    options?: { status?: number; code?: string },
  ) {
    super(message)
    this.name = 'QrApiError'
    this.kind = kind
    this.status = options?.status
    this.code = options?.code
  }
}

async function parseJsonSafely(response: Response): Promise<unknown | null> {
  try {
    const text = await response.text()
    if (!text) {
      return null
    }

    return JSON.parse(text) as unknown
  } catch {
    return null
  }
}

export async function analyzeQrMatrix(request: QrRequest): Promise<QrResponse> {
  let response: Response

  try {
    response = await fetch(`${apiBaseUrl}/api/v1/qr`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    })
  } catch {
    throw new QrApiError(
      'network',
      'Unable to reach the QR service. Check your connection and try again.',
    )
  }

  const payload = await parseJsonSafely(response)

  if (!response.ok) {
    if (isApiErrorEnvelope(payload)) {
      throw new QrApiError('http', payload.error.message, {
        status: response.status,
        code: payload.error.code,
      })
    }

    throw new QrApiError(
      'http',
      `Request failed with status ${response.status}. Please try again.`,
      { status: response.status },
    )
  }

  if (!isQrResponse(payload)) {
    throw new QrApiError(
      'invalid-response',
      'The QR service returned an unexpected response. Please try again.',
      { status: response.status },
    )
  }

  return payload
}
