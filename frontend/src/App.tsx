import { useMemo, useState } from 'react'
import { MatrixEditor } from './components/MatrixEditor'

const MIN_DIMENSION = 1
const MAX_DIMENSION = 10
const DEFAULT_ROWS = 2
const DEFAULT_COLUMNS = 2

type MatrixInput = string[][]

function createMatrix(rows: number, columns: number, defaultValue = '0'): MatrixInput {
  return Array.from({ length: rows }, () =>
    Array.from({ length: columns }, () => defaultValue),
  )
}

function clampDimension(value: number): number {
  if (!Number.isFinite(value)) {
    return MIN_DIMENSION
  }

  return Math.max(MIN_DIMENSION, Math.min(MAX_DIMENSION, Math.trunc(value)))
}

function resizeMatrix(
  matrix: MatrixInput,
  nextRows: number,
  nextColumns: number,
  defaultValue = '0',
): MatrixInput {
  return Array.from({ length: nextRows }, (_, rowIndex) =>
    Array.from({ length: nextColumns }, (_, columnIndex) => {
      const existingRow = matrix[rowIndex]
      if (existingRow && typeof existingRow[columnIndex] === 'string') {
        return existingRow[columnIndex]
      }
      return defaultValue
    }),
  )
}

function isCellValueValid(rawValue: string): boolean {
  const trimmed = rawValue.trim()
  if (trimmed === '') {
    return false
  }

  return Number.isFinite(Number(trimmed))
}

function App() {
  const [rows, setRows] = useState<number>(DEFAULT_ROWS)
  const [columns, setColumns] = useState<number>(DEFAULT_COLUMNS)
  const [rowsInput, setRowsInput] = useState<string>(String(DEFAULT_ROWS))
  const [columnsInput, setColumnsInput] = useState<string>(String(DEFAULT_COLUMNS))
  const [matrix, setMatrix] = useState<MatrixInput>(
    createMatrix(DEFAULT_ROWS, DEFAULT_COLUMNS),
  )
  const hasInvalidCells = matrix.some((row) =>
    row.some((value) => !isCellValueValid(value)),
  )

  const dimensionLimits = useMemo(
    () => ({ min: MIN_DIMENSION, max: MAX_DIMENSION }),
    [],
  )

  const updateDimensions = (nextRowsInput: number, nextColumnsInput: number) => {
    const nextRows = clampDimension(nextRowsInput)
    const nextColumns = clampDimension(nextColumnsInput)

    setRows(nextRows)
    setColumns(nextColumns)
    setRowsInput(String(nextRows))
    setColumnsInput(String(nextColumns))
    setMatrix((current) => resizeMatrix(current, nextRows, nextColumns))
  }

  const handleRowsCommit = (nextRows: number) => {
    updateDimensions(nextRows, columns)
  }

  const handleColumnsCommit = (nextColumns: number) => {
    updateDimensions(rows, nextColumns)
  }

  const handleRowsInputChange = (nextValue: string) => {
    setRowsInput(nextValue)
  }

  const handleColumnsInputChange = (nextValue: string) => {
    setColumnsInput(nextValue)
  }

  const handleRowsInputCommit = () => {
    const trimmed = rowsInput.trim()
    if (trimmed === '') {
      setRowsInput(String(rows))
      return
    }

    const parsed = Number(trimmed)
    if (!Number.isFinite(parsed)) {
      setRowsInput(String(rows))
      return
    }

    handleRowsCommit(parsed)
  }

  const handleColumnsInputCommit = () => {
    const trimmed = columnsInput.trim()
    if (trimmed === '') {
      setColumnsInput(String(columns))
      return
    }

    const parsed = Number(trimmed)
    if (!Number.isFinite(parsed)) {
      setColumnsInput(String(columns))
      return
    }

    handleColumnsCommit(parsed)
  }

  const handleCellChange = (
    rowIndex: number,
    columnIndex: number,
    rawValue: string,
  ) => {
    setMatrix((current) => {
      const next = current.map((row) => [...row])
      next[rowIndex][columnIndex] = rawValue
      return next
    })
  }

  const handleLoadSample = () => {
    const sample: MatrixInput = [
      ['1', '2', '3'],
      ['4', '5', '6'],
    ]

    setRows(sample.length)
    setColumns(sample[0].length)
    setRowsInput(String(sample.length))
    setColumnsInput(String(sample[0].length))
    setMatrix(sample.map((row) => [...row]))
  }

  const handleReset = () => {
    setRows(DEFAULT_ROWS)
    setColumns(DEFAULT_COLUMNS)
    setRowsInput(String(DEFAULT_ROWS))
    setColumnsInput(String(DEFAULT_COLUMNS))
    setMatrix(createMatrix(DEFAULT_ROWS, DEFAULT_COLUMNS))
  }

  return (
    <main className="app-shell">
      <header className="app-header">
        <h1>Matrix QR Analytics</h1>
        <p>Analyze rectangular matrices using QR decomposition.</p>
      </header>

      <MatrixEditor
        rowsInput={rowsInput}
        columnsInput={columnsInput}
        matrix={matrix}
        dimensionLimits={dimensionLimits}
        hasInvalidCells={hasInvalidCells}
        isCellValueValid={isCellValueValid}
        onRowsInputChange={handleRowsInputChange}
        onColumnsInputChange={handleColumnsInputChange}
        onRowsInputCommit={handleRowsInputCommit}
        onColumnsInputCommit={handleColumnsInputCommit}
        onCellChange={handleCellChange}
        onLoadSample={handleLoadSample}
        onReset={handleReset}
      />
    </main>
  )
}

export default App
