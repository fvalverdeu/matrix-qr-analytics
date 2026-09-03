import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import { analyzeQrMatrix, QrApiError } from './api/qrApi'

vi.mock('./api/qrApi', () => {
  class MockQrApiError extends Error {
    readonly kind: 'network' | 'http' | 'invalid-response'
    readonly status?: number
    readonly code?: string

    constructor(
      kind: 'network' | 'http' | 'invalid-response',
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

  return {
    analyzeQrMatrix: vi.fn(),
    QrApiError: MockQrApiError,
  }
})

type Deferred<T> = {
  promise: Promise<T>
  resolve: (value: T) => void
  reject: (reason?: unknown) => void
}

function createDeferred<T>(): Deferred<T> {
  let resolve!: (value: T) => void
  let reject!: (reason?: unknown) => void

  const promise = new Promise<T>((res, rej) => {
    resolve = res
    reject = rej
  })

  return { promise, resolve, reject }
}

describe('App', () => {
  const analyzeQrMatrixMock = vi.mocked(analyzeQrMatrix)

  beforeEach(() => {
    analyzeQrMatrixMock.mockReset()
  })

  afterEach(() => {
    vi.clearAllMocks()
    cleanup()
  })

  it('starts with a 2x2 matrix and enabled Analyze button', () => {
    render(<App />)

    const cellInputs = screen.getAllByRole('textbox')
    expect(cellInputs).toHaveLength(4)

    const analyzeButton = screen.getByRole('button', { name: /analyze matrix/i })
    expect(analyzeButton).toBeEnabled()
  })

  it('disables Analyze when a cell is cleared and shows validation message', () => {
    render(<App />)

    const firstCell = screen.getByLabelText('Cell row 1, column 1')
    fireEvent.change(firstCell, { target: { value: '' } })

    const analyzeButton = screen.getByRole('button', { name: /analyze matrix/i })
    expect(analyzeButton).toBeDisabled()

    expect(
      screen.getByText('Please enter finite numeric values in all cells before analysis.'),
    ).toBeInTheDocument()
  })

  it('loads a 2x3 sample matrix with expected values', () => {
    render(<App />)

    fireEvent.click(screen.getByRole('button', { name: /load sample matrix/i }))

    const cellInputs = screen.getAllByRole('textbox')
    expect(cellInputs).toHaveLength(6)
    expect(screen.getByLabelText('Cell row 1, column 1')).toHaveValue('1')
    expect(screen.getByLabelText('Cell row 1, column 2')).toHaveValue('2')
    expect(screen.getByLabelText('Cell row 1, column 3')).toHaveValue('3')
    expect(screen.getByLabelText('Cell row 2, column 1')).toHaveValue('4')
    expect(screen.getByLabelText('Cell row 2, column 2')).toHaveValue('5')
    expect(screen.getByLabelText('Cell row 2, column 3')).toHaveValue('6')
  })

  it('shows loading state and then renders success results', async () => {
    const successfulResponse = {
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

    const deferred = createDeferred<typeof successfulResponse>()

    analyzeQrMatrixMock.mockReturnValueOnce(deferred.promise)

    render(<App />)

    fireEvent.click(screen.getByRole('button', { name: /analyze matrix/i }))

    const loadingButton = screen.getByRole('button', { name: /analyzing/i })
    expect(loadingButton).toBeDisabled()

    deferred.resolve(successfulResponse)

    await waitFor(() => {
      expect(screen.getByText('Analysis completed successfully.')).toBeInTheDocument()
    })

    expect(screen.getByRole('heading', { name: 'Q' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'R' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Statistics' })).toBeInTheDocument()
  })

  it('shows user-facing error and no results on API failure', async () => {
    analyzeQrMatrixMock.mockRejectedValueOnce(
      new QrApiError('http', 'Matrix cannot be empty', {
        status: 400,
        code: 'INVALID_MATRIX',
      }),
    )

    render(<App />)

    fireEvent.click(screen.getByRole('button', { name: /analyze matrix/i }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Matrix cannot be empty')
    })

    expect(
      screen.queryByRole('heading', { name: 'Analysis results' }),
    ).not.toBeInTheDocument()
  })

  it('clears stale results when matrix input changes', async () => {
    analyzeQrMatrixMock.mockResolvedValueOnce({
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
    })

    render(<App />)

    fireEvent.click(screen.getByRole('button', { name: /analyze matrix/i }))

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Analysis results' })).toBeInTheDocument()
    })

    fireEvent.change(screen.getByLabelText('Cell row 1, column 1'), {
      target: { value: '9' },
    })

    expect(
      screen.queryByRole('heading', { name: 'Analysis results' }),
    ).not.toBeInTheDocument()
    expect(screen.queryByText('Analysis completed successfully.')).not.toBeInTheDocument()
  })
})
