const test = require('node:test');
const assert = require('node:assert/strict');

const { createStatisticsController } = require('./statistics.controller');

test('calculateStatistics maps unexpected service errors to INTERNAL_ERROR without leaking internals', () => {
  const controller = createStatisticsController({
    calculateStatistics() {
      throw new Error('sensitive internal failure');
    },
  });

  const req = {
    body: {
      q: [[1]],
      r: [[1]],
    },
  };

  const res = {
    statusCode: 200,
    payload: undefined,
    status(code) {
      this.statusCode = code;
      return this;
    },
    json(body) {
      this.payload = body;
      return this;
    },
  };

  controller.calculateStatistics(req, res);

  assert.equal(res.statusCode, 500);
  assert.deepEqual(res.payload, {
    error: {
      code: 'INTERNAL_ERROR',
      message: 'An unexpected error occurred while calculating statistics',
    },
  });

  assert.equal(
    JSON.stringify(res.payload).includes('sensitive internal failure'),
    false,
  );
});
