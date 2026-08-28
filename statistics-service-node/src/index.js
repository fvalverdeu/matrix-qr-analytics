require('dotenv').config();

const express = require('express');
const config = require('../config');
const { createStatisticsController } = require('./controllers/statistics.controller');
const { createStatisticsRoutes } = require('./routes/statistics.routes');
const statisticsService = require('./services/statistics.service');

const app = express();

app.use(express.json());

app.use((err, req, res, next) => {
  if (err instanceof SyntaxError && err.status === 400 && 'body' in err) {
    return res.status(400).json({
      error: {
        code: 'INVALID_REQUEST',
        message: 'Request body must be valid JSON',
      },
    });
  }

  return next(err);
});

const statisticsController = createStatisticsController(statisticsService);
app.use(createStatisticsRoutes(statisticsController));

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'statistics-service-node' });
});

app.listen(config.port, () => {
  console.log(`Statistics service listening on port ${config.port}`);
});
