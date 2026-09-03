import type { QrResponse } from '../types/qr'
import { ResultMatrix } from './ResultMatrix'
import { StatisticsPanel } from './StatisticsPanel'

type AnalysisResultsProps = {
  result: QrResponse
}

export function AnalysisResults({ result }: AnalysisResultsProps) {
  return (
    <section className="results" aria-labelledby="analysis-results-heading">
      <h2 id="analysis-results-heading">Analysis results</h2>
      <div className="results-grid">
        <ResultMatrix title="Q" matrix={result.q} />
        <ResultMatrix title="R" matrix={result.r} />
      </div>
      <StatisticsPanel statistics={result.statistics} />
    </section>
  )
}
