import type { QrStatistics } from '../types/qr'
import { formatNumberForDisplay } from '../utils/numberFormat'

type StatisticsPanelProps = {
  statistics: QrStatistics
}

export function StatisticsPanel({ statistics }: StatisticsPanelProps) {
  return (
    <section className="result-card" aria-label="Statistics">
      <h3>Statistics</h3>
      <dl className="statistics-grid">
        <dt>Maximum</dt>
        <dd>{formatNumberForDisplay(statistics.max)}</dd>

        <dt>Minimum</dt>
        <dd>{formatNumberForDisplay(statistics.min)}</dd>

        <dt>Average</dt>
        <dd>{formatNumberForDisplay(statistics.average)}</dd>

        <dt>Sum</dt>
        <dd>{formatNumberForDisplay(statistics.sum)}</dd>

        <dt>Has diagonal matrix</dt>
        <dd>{statistics.hasDiagonalMatrix ? 'Yes' : 'No'}</dd>
      </dl>
    </section>
  )
}
