import express, { type Request, type Response } from 'express';

interface StatisticsController {
  calculateStatistics(req: Request, res: Response): Response;
}

export function createStatisticsRoutes(statisticsController: StatisticsController) {
  const router = express.Router();

  router.post('/api/v1/statistics', statisticsController.calculateStatistics);

  return router;
}
