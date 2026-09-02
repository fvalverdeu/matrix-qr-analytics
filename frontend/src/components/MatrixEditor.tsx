import { MatrixGrid } from './MatrixGrid'

type MatrixEditorProps = {
  rowsInput: string
  columnsInput: string
  matrix: string[][]
  dimensionLimits: {
    min: number
    max: number
  }
  hasInvalidCells: boolean
  isCellValueValid: (value: string) => boolean
  onRowsInputChange: (value: string) => void
  onColumnsInputChange: (value: string) => void
  onRowsInputCommit: () => void
  onColumnsInputCommit: () => void
  onCellChange: (rowIndex: number, columnIndex: number, rawValue: string) => void
  onLoadSample: () => void
  onReset: () => void
}

export function MatrixEditor({
  rowsInput,
  columnsInput,
  matrix,
  dimensionLimits,
  hasInvalidCells,
  isCellValueValid,
  onRowsInputChange,
  onColumnsInputChange,
  onRowsInputCommit,
  onColumnsInputCommit,
  onCellChange,
  onLoadSample,
  onReset,
}: MatrixEditorProps) {
  return (
    <section className="editor" aria-label="Matrix editor">
      <div className="editor-controls">
        <div className="dimension-control">
          <label htmlFor="matrix-rows">Rows</label>
          <input
            id="matrix-rows"
            type="number"
            min={dimensionLimits.min}
            max={dimensionLimits.max}
            step={1}
            value={rowsInput}
            onChange={(event) => onRowsInputChange(event.target.value)}
            onBlur={onRowsInputCommit}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                onRowsInputCommit()
              }
            }}
          />
        </div>

        <div className="dimension-control">
          <label htmlFor="matrix-columns">Columns</label>
          <input
            id="matrix-columns"
            type="number"
            min={dimensionLimits.min}
            max={dimensionLimits.max}
            step={1}
            value={columnsInput}
            onChange={(event) => onColumnsInputChange(event.target.value)}
            onBlur={onColumnsInputCommit}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                onColumnsInputCommit()
              }
            }}
          />
        </div>

        <div className="editor-actions">
          <button type="button" onClick={onLoadSample}>
            Load sample matrix
          </button>
          <button type="button" onClick={onReset}>
            Reset to 2x2
          </button>
        </div>
      </div>

      <MatrixGrid
        matrix={matrix}
        isCellValueValid={isCellValueValid}
        onCellChange={onCellChange}
      />

      <div className="editor-footer">
        <p className="editor-hint" role="status" aria-live="polite">
          {hasInvalidCells
            ? 'Please enter finite numeric values in all cells before analysis.'
            : 'Matrix input is ready. API analysis will be enabled in the integration step.'}
        </p>

        <button type="button" className="primary-action" disabled>
          Analyze matrix (available in integration step)
        </button>
      </div>
    </section>
  )
}
