const { validateMatrices } = require('./matrix.validator');

const DIAGONAL_EPSILON = 1e-10;

function isDiagonal(matrix) {
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

function collectValues(q, r) {
  const values = [];

  for (const row of q) {
    values.push(...row);
  }

  for (const row of r) {
    values.push(...row);
  }

  return values;
}

function calculateStatistics(q, r) {
  validateMatrices(q, r);

  const values = collectValues(q, r);
  const sum = values.reduce((total, value) => total + value, 0);

  return {
    max: Math.max(...values),
    min: Math.min(...values),
    average: sum / values.length,
    sum,
    hasDiagonalMatrix: isDiagonal(q) || isDiagonal(r),
  };
}

module.exports = { calculateStatistics };
