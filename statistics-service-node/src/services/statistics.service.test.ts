import assert from 'node:assert/strict';
import test from 'node:test';

import { ValidationError } from '../errors/validation.error';
import { calculateStatistics } from './statistics.service';

function assertAlmostEqual(actual: number, expected: number, tolerance = 1e-12) {
  assert.ok(Math.abs(actual - expected) <= tolerance, `expected ${actual} to be within ${tolerance} of ${expected}`);
}

test('calculateStatistics handles scalar matrices', () => {
  const result = calculateStatistics([[2]], [[3]]);

  assert.equal(result.max, 3);
  assert.equal(result.min, 2);
  assert.equal(result.sum, 5);
  assert.equal(result.average, 2.5);
  assert.equal(result.hasDiagonalMatrix, true);
});

test('calculateStatistics returns expected contract fields', () => {
  const result = calculateStatistics([[1]], [[2]]);
  const keys = Object.keys(result).sort();

  assert.deepEqual(keys, ['average', 'hasDiagonalMatrix', 'max', 'min', 'sum']);
});

test('calculateStatistics computes statistics over mixed square matrices', () => {
  const q = [
    [1, 2],
    [3, 4],
  ];
  const r = [
    [-5, 6],
    [7, 8],
  ];

  const result = calculateStatistics(q, r);

  assert.equal(result.max, 8);
  assert.equal(result.min, -5);
  assert.equal(result.sum, 26);
  assert.equal(result.average, 3.25);
  assert.equal(result.hasDiagonalMatrix, false);
});

test('calculateStatistics handles negative and decimal values', () => {
  const q = [
    [-1.5, 2.25],
    [3.75, -4.5],
  ];
  const r = [
    [0.5, -0.25],
    [1.25, 2.75],
  ];

  const result = calculateStatistics(q, r);

  assert.equal(result.max, 3.75);
  assert.equal(result.min, -4.5);
  assert.equal(result.sum, 4.25);
  assert.equal(result.average, 0.53125);
  assert.equal(result.hasDiagonalMatrix, false);
});

test('calculateStatistics handles zero values', () => {
  const q = [
    [0, 0],
    [0, 0],
  ];
  const r = [[0]];

  const result = calculateStatistics(q, r);

  assert.equal(result.max, 0);
  assert.equal(result.min, 0);
  assert.equal(result.sum, 0);
  assert.equal(result.average, 0);
  assert.equal(result.hasDiagonalMatrix, true);
});

test('calculateStatistics supports differently shaped rectangular matrices', () => {
  const q = [[1, 2, 3]];
  const r = [
    [4, 5],
    [6, 7],
  ];

  const result = calculateStatistics(q, r);

  assert.equal(result.max, 7);
  assert.equal(result.min, 1);
  assert.equal(result.sum, 28);
  assert.equal(result.average, 4);
  assert.equal(result.hasDiagonalMatrix, false);
});

test('calculateStatistics aggregates values from both q and r', () => {
  const q = [[-10, 2]];
  const r = [[3, 99]];

  const result = calculateStatistics(q, r);

  assert.equal(result.min, -10);
  assert.equal(result.max, 99);
  assert.equal(result.sum, 94);
  assert.equal(result.average, 23.5);
});

test('calculateStatistics diagonal semantics are preserved', () => {
  const testCases = [
    {
      name: 'q diagonal, r not diagonal',
      q: [
        [1, 0],
        [0, 2],
      ],
      r: [
        [1, 3],
        [0, 1],
      ],
      expected: true,
    },
    {
      name: 'q not diagonal, r diagonal',
      q: [
        [1, 3],
        [0, 1],
      ],
      r: [
        [5, 0],
        [0, 6],
      ],
      expected: true,
    },
    {
      name: 'neither diagonal',
      q: [
        [1, 3],
        [0, 2],
      ],
      r: [
        [4, 5],
        [6, 7],
      ],
      expected: false,
    },
    {
      name: 'both diagonal',
      q: [
        [1, 0],
        [0, 2],
      ],
      r: [
        [3, 0],
        [0, 4],
      ],
      expected: true,
    },
    {
      name: 'non-square with zero off conceptual diagonal is not diagonal',
      q: [[1, 0, 0]],
      r: [[2, 3, 4]],
      expected: false,
    },
    {
      name: 'off-diagonal exactly zero is diagonal when square',
      q: [
        [1, 0],
        [0, 2],
      ],
      r: [[3, 4, 5]],
      expected: true,
    },
    {
      name: 'off-diagonal below epsilon is diagonal',
      q: [
        [1, 1e-12],
        [0, 2],
      ],
      r: [[3, 4, 5]],
      expected: true,
    },
    {
      name: 'off-diagonal above epsilon is not diagonal',
      q: [
        [1, 1e-5],
        [0, 2],
      ],
      r: [[3, 4, 5]],
      expected: false,
    },
    {
      name: 'off-diagonal at epsilon boundary is diagonal',
      q: [
        [1, 1e-10],
        [0, 2],
      ],
      r: [[3, 4, 5]],
      expected: true,
    },
  ];

  for (const tc of testCases) {
    const result = calculateStatistics(tc.q, tc.r);
    assert.equal(result.hasDiagonalMatrix, tc.expected, tc.name);
  }
});

test('calculateStatistics propagates ValidationError for invalid q', () => {
  assert.throws(
    () => calculateStatistics(null, [[1]]),
    (error: unknown) => {
      assert.ok(error instanceof ValidationError, `expected ValidationError, got ${(error as { constructor?: { name?: string } })?.constructor?.name}`);
      assert.equal(error.message, 'q is required');
      return true;
    }
  );
});

test('calculateStatistics propagates ValidationError for invalid r', () => {
  assert.throws(
    () => calculateStatistics([[1]], []),
    (error: unknown) => {
      assert.ok(error instanceof ValidationError, `expected ValidationError, got ${(error as { constructor?: { name?: string } })?.constructor?.name}`);
      assert.equal(error.message, 'r must be a valid matrix');
      return true;
    }
  );
});
