import express, { type ErrorRequestHandler } from 'express';

import { createStatisticsController } from './controllers/statistics.controller';
import { createStatisticsRoutes } from './routes/statistics.routes';
import * as statisticsService from './services/statistics.service';

function isMalformedJsonError(error: unknown): error is SyntaxError & { status: number; body: unknown } {
  if (!(error instanceof SyntaxError)) {
    return false;
  }

  if (typeof error !== 'object' || error === null) {
    return false;
  }

  if (!('status' in error) || error.status !== 400) {
    return false;
  }

  return 'body' in error;
}

export function createApp() {
  const app = express();

  app.use(express.json());

  const malformedJsonHandler: ErrorRequestHandler = (error, req, res, next) => {
    if (isMalformedJsonError(error)) {
      return res.status(400).json({
        error: {
          code: 'INVALID_REQUEST',
          message: 'Request body must be valid JSON',
        },
      });
    }

    return next(error);
  };

  app.use(malformedJsonHandler);

  const statisticsController = createStatisticsController(statisticsService);
  app.use(createStatisticsRoutes(statisticsController));

  app.get('/health', (req, res) => {
    res.json({ status: 'ok', service: 'statistics-service-node' });
  });

  return app;
}
