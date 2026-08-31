import assert from 'node:assert/strict';
import test from 'node:test';

import { ValidationError } from '../errors/validation.error';
import { validateMatrices, validateMatrix } from './matrix.validator';

function assertValidationError(error: unknown, expectedMessage: string) {
  assert.ok(error instanceof ValidationError, `expected ValidationError, got ${(error as { constructor?: { name?: string } })?.constructor?.name}`);
  assert.equal(error.message, expectedMessage);
}

test('validateMatrix accepts representative valid matrices', () => {
  const testCases = [
    {
      name: 'scalar_1x1',
      matrix: [[5]],
    },
    {
      name: 'square_2x2',
      matrix: [
        [1, 2],
        [3, 4],
      ],
    },
    {
      name: 'tall_3x2',
      matrix: [
        [1, 2],
        [3, 4],
        [5, 6],
      ],
    },
    {
      name: 'wide_2x3',
      matrix: [
        [1, 2, 3],
        [4, 5, 6],
      ],
    },
    {
      name: 'zero_matrix',
      matrix: [
        [0, 0],
        [0, 0],
      ],
    },
    {
      name: 'negative_decimal_values',
      matrix: [
        [-1.5, 2.25],
        [3.75, -4.5],
      ],
    },
  ];

  for (const tc of testCases) {
    assert.doesNotThrow(() => validateMatrix(tc.matrix, 'q'), tc.name);
  }
});

test('validateMatrix required-field messages for q/r', () => {
  assert.throws(
    () => validateMatrix(undefined, 'q'),
    (error: unknown) => {
      assertValidationError(error, 'q is required');
      return true;
    }
  );

  assert.throws(
    () => validateMatrix(null, 'r'),
    (error: unknown) => {
      assertValidationError(error, 'r is required');
      return true;
    }
  );
});

test('validateMatrix rejects intended invalid matrices with ValidationError', () => {
  const testCases = [
    {
      name: 'null_matrix',
      matrix: null,
      nameArg: 'q',
      message: 'q is required',
    },
    {
      name: 'empty_matrix',
      matrix: [],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'empty_first_row',
      matrix: [[]],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'empty_later_row',
      matrix: [[1], []],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'ragged_rows',
      matrix: [
        [1, 2],
        [3],
      ],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'nonnumeric_cell',
      matrix: [[1, '2']],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'null_cell',
      matrix: [[1, null]],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'boolean_cell',
      matrix: [[1, true]],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'invalid_later_row_null',
      matrix: [
        [1, 2],
        null,
      ],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
    {
      name: 'invalid_later_row_string',
      matrix: [
        [1, 2],
        'bad',
      ],
      nameArg: 'q',
      message: 'q must be a valid matrix',
    },
  ];

  for (const tc of testCases) {
    assert.throws(
      () => validateMatrix(tc.matrix, tc.nameArg),
      (error: unknown) => {
        assertValidationError(error, tc.message);
        return true;
      },
      tc.name
    );
  }
});

test('validateMatrices validates q and r orchestration', () => {
  assert.doesNotThrow(() => {
    validateMatrices(
      [
        [1, 2],
        [3, 4],
      ],
      [
        [5, 6],
        [7, 8],
      ]
    );
  });

  assert.throws(
    () => validateMatrices([], [[1]]),
    (error: unknown) => {
      assertValidationError(error, 'q must be a valid matrix');
      return true;
    }
  );

  assert.throws(
    () => validateMatrices([[1]], []),
    (error: unknown) => {
      assertValidationError(error, 'r must be a valid matrix');
      return true;
    }
  );
});

test('first row null should produce ValidationError', () => {
  assert.throws(
    () => validateMatrix([null], 'q'),
    (error: unknown) => {
      assertValidationError(error, 'q must be a valid matrix');
      return true;
    }
  );
});
