const { ValidationError } = require('../errors/validation.error');

function validateMatrix(matrix, name) {
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

function validateMatrices(q, r) {
  validateMatrix(q, 'q');
  validateMatrix(r, 'r');
}

module.exports = { validateMatrix, validateMatrices };
