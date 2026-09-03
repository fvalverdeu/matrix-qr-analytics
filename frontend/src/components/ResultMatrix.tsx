import { formatNumberForDisplay } from '../utils/numberFormat'

type ResultMatrixProps = {
  title: string
  matrix: number[][]
}

export function ResultMatrix({ title, matrix }: ResultMatrixProps) {
  return (
    <section className="result-card" aria-label={`${title} matrix`}>
      <h3>{title}</h3>
      <div className="result-matrix-wrapper" role="region" aria-label={`${title} matrix values`}>
        <table className="result-matrix-table">
          <caption className="sr-only">{title} matrix result</caption>
          <tbody>
            {matrix.map((row, rowIndex) => (
              <tr key={`${title}-row-${rowIndex}`}>
                {row.map((value, columnIndex) => (
                  <td key={`${title}-cell-${rowIndex}-${columnIndex}`}>
                    {formatNumberForDisplay(value)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}
