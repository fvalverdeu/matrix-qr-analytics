type MatrixGridProps = {
  matrix: string[][]
  isCellValueValid: (value: string) => boolean
  onCellChange: (rowIndex: number, columnIndex: number, rawValue: string) => void
}

export function MatrixGrid({
  matrix,
  isCellValueValid,
  onCellChange,
}: MatrixGridProps) {
  return (
    <div className="matrix-grid-wrapper" role="region" aria-label="Matrix input grid">
      <table className="matrix-grid">
        <caption className="sr-only">Editable matrix values</caption>
        <tbody>
          {matrix.map((row, rowIndex) => (
            <tr key={`row-${rowIndex}`}>
              {row.map((value, columnIndex) => {
                const hasError = !isCellValueValid(value)
                const inputId = `matrix-cell-${rowIndex}-${columnIndex}`

                return (
                  <td key={inputId}>
                    <label htmlFor={inputId} className="sr-only">
                      Cell row {rowIndex + 1}, column {columnIndex + 1}
                    </label>
                    <input
                      id={inputId}
                      type="text"
                      inputMode="decimal"
                      value={value}
                      aria-invalid={hasError}
                      className={hasError ? 'cell-input cell-input-invalid' : 'cell-input'}
                      onChange={(event) =>
                        onCellChange(rowIndex, columnIndex, event.target.value)
                      }
                    />
                  </td>
                )
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
