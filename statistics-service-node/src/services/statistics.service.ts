import { validateMatrices } from './matrix.validator';
import type { Matrix, Statistics } from '../types';

const DIAGONAL_EPSILON = 1e-10;

function isDiagonal(matrix: Matrix): boolean {
  const rows = matrix.length;
  const cols = matrix[0].length;

  if (rows !== cols) {
    return false;
  }

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > DIAGONAL_EPSILON) {
        return false;
      }
    }
  }

  return true;
}

function collectValues(q: Matrix, r: Matrix): number[] {
  const values: number[] = [];

  for (const row of q) {
    values.push(...row);
  }

  for (const row of r) {
    values.push(...row);
  }

  return values;
}

export function calculateStatistics(q: unknown, r: unknown): Statistics {
  const validated = validateMatrices(q, r);

  const values = collectValues(validated.q, validated.r);
  const sum = values.reduce((total, value) => total + value, 0);

  return {
    max: Math.max(...values),
    min: Math.min(...values),
    average: sum / values.length,
    sum,
    hasDiagonalMatrix: isDiagonal(validated.q) || isDiagonal(validated.r),
  };
}
