import type { Request, Response } from 'express';

import type { Statistics } from '../types';
import { ValidationError } from '../errors/validation.error';

interface StatisticsService {
  calculateStatistics(q: unknown, r: unknown): Statistics;
}

function extractMatricesFromBody(body: unknown): { q: unknown; r: unknown } {
  if (body === null || typeof body !== 'object') {
    return { q: undefined, r: undefined };
  }

  const record = body as Record<string, unknown>;
  return {
    q: record.q,
    r: record.r,
  };
}

export function createStatisticsController(statisticsService: StatisticsService) {
  function calculateStatistics(req: Request, res: Response): Response {
    const { q, r } = extractMatricesFromBody(req.body as unknown);

    try {
      const statistics = statisticsService.calculateStatistics(q, r);
      return res.json(statistics);
    } catch (error: unknown) {
      if (error instanceof ValidationError) {
        return res.status(400).json({
          error: {
            code: 'INVALID_MATRICES',
            message: error.message,
          },
        });
      }

      return res.status(500).json({
        error: {
          code: 'INTERNAL_ERROR',
          message: 'An unexpected error occurred while calculating statistics',
        },
      });
    }
  }

  return { calculateStatistics };
}
