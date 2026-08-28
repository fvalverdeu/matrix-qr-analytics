const { ValidationError } = require('../errors/validation.error');

function createStatisticsController(statisticsService) {
  function calculateStatistics(req, res) {
    const { q, r } = req.body ?? {};

    try {
      const statistics = statisticsService.calculateStatistics(q, r);
      return res.json(statistics);
    } catch (error) {
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

module.exports = { createStatisticsController };
