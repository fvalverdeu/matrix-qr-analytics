import assert from 'node:assert/strict';
import test from 'node:test';

import type { Request, Response } from 'express';

import { createStatisticsController } from './statistics.controller';

interface MockResponse {
  statusCode: number;
  payload: unknown;
  status(this: MockResponse, code: number): MockResponse;
  json(this: MockResponse, body: unknown): MockResponse;
}

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
  } as Request;

  const res: MockResponse = {
    statusCode: 200,
    payload: undefined,
    status(code: number) {
      this.statusCode = code;
      return this;
    },
    json(body: unknown) {
      this.payload = body;
      return this;
    },
  };

  controller.calculateStatistics(req, res as unknown as Response);

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
