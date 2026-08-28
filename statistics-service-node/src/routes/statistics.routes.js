const express = require('express');

function createStatisticsRoutes(statisticsController) {
  const router = express.Router();

  router.post('/api/v1/statistics', statisticsController.calculateStatistics);

  return router;
}

module.exports = { createStatisticsRoutes };
