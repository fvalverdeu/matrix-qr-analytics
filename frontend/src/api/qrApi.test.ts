import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { apiBaseUrl } from '../config/env'
import { analyzeQrMatrix, QrApiError } from './qrApi'

const validResponse = {
  q: [
    [1, 0],
    [0, 1],
  ],
  r: [
    [1, 2],
    [0, 3],
  ],
  statistics: {
    max: 3,
    min: 0,
    average: 1.166667,
    sum: 7,
    hasDiagonalMatrix: true,
  },
}

describe('analyzeQrMatrix', () => {
  const fetchMock = vi.fn<typeof fetch>()

  beforeEach(() => {
    fetchMock.mockReset()
    vi.stubGlobal('fetch', fetchMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('sends POST request with JSON body and returns valid response', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify(validResponse), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    const request = { matrix: [[1, 2], [3, 4]] }
    const result = await analyzeQrMatrix(request)

    expect(result).toEqual(validResponse)
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock).toHaveBeenCalledWith(`${apiBaseUrl}/api/v1/qr`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    })
  })

  it('returns http error preserving backend envelope details', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          error: {
            code: 'INVALID_MATRIX',
            message: 'Matrix cannot be empty',
          },
        }),
        {
          status: 400,
          headers: { 'Content-Type': 'application/json' },
        },
      ),
    )

    await expect(analyzeQrMatrix({ matrix: [[1]] })).rejects.toMatchObject({
      name: 'QrApiError',
      kind: 'http',
      status: 400,
      code: 'INVALID_MATRIX',
      message: 'Matrix cannot be empty',
    })
  })

  it('returns safe http fallback for malformed response body', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response('not-json', {
        status: 502,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(analyzeQrMatrix({ matrix: [[1]] })).rejects.toMatchObject({
      name: 'QrApiError',
      kind: 'http',
      status: 502,
      message: 'Request failed with status 502. Please try again.',
    })
  })

  it('returns network error for fetch failures', async () => {
    fetchMock.mockRejectedValueOnce(new Error('network down'))

    await expect(analyzeQrMatrix({ matrix: [[1]] })).rejects.toMatchObject({
      name: 'QrApiError',
      kind: 'network',
      message: 'Unable to reach the QR service. Check your connection and try again.',
    })
  })

  it('returns invalid-response error for successful response with invalid contract', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ ok: true }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(analyzeQrMatrix({ matrix: [[1]] })).rejects.toMatchObject({
      name: 'QrApiError',
      kind: 'invalid-response',
      status: 200,
      message: 'The QR service returned an unexpected response. Please try again.',
    })
  })

  it('returns safe http fallback for empty error body', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response('', {
        status: 500,
      }),
    )

    try {
      await analyzeQrMatrix({ matrix: [[1]] })
      throw new Error('expected analyzeQrMatrix to fail')
    } catch (error) {
      expect(error).toBeInstanceOf(QrApiError)
      expect((error as QrApiError).kind).toBe('http')
      expect((error as QrApiError).status).toBe(500)
      expect((error as QrApiError).message).toBe(
        'Request failed with status 500. Please try again.',
      )
    }
  })
})
