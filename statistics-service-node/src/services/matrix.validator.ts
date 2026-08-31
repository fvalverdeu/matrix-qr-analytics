import { ValidationError } from '../errors/validation.error';
import type { Matrix } from '../types';

export function validateMatrix(matrix: unknown, name: string): asserts matrix is Matrix {
  if (matrix === undefined || matrix === null) {
    throw new ValidationError(`${name} is required`);
  }

  if (!Array.isArray(matrix)) {
    throw new ValidationError(`${name} must be a valid matrix`);
  }

  if (matrix.length === 0) {
    throw new ValidationError(`${name} must be a valid matrix`);
  }

  const firstRow = matrix[0];
  if (!Array.isArray(firstRow) || firstRow.length === 0) {
    throw new ValidationError(`${name} must be a valid matrix`);
  }

  const cols = firstRow.length;

  for (const row of matrix) {
    if (!Array.isArray(row) || row.length === 0) {
      throw new ValidationError(`${name} must be a valid matrix`);
    }

    if (row.length !== cols) {
      throw new ValidationError(`${name} must be a valid matrix`);
    }

    for (const value of row) {
      if (typeof value !== 'number' || !Number.isFinite(value)) {
        throw new ValidationError(`${name} must be a valid matrix`);
      }
    }
  }
}

export function validateMatrices(q: unknown, r: unknown): { q: Matrix; r: Matrix } {
  validateMatrix(q, 'q');
  validateMatrix(r, 'r');

  return { q, r };
}
