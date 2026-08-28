package services

import "math"

type MatrixValidator struct{}

func NewMatrixValidator() *MatrixValidator {
	return &MatrixValidator{}
}

func (v *MatrixValidator) Validate(matrix [][]float64) error {
	if matrix == nil {
		return newValidationError("Matrix is required")
	}

	if len(matrix) == 0 {
		return newValidationError("Matrix cannot be empty")
	}

	cols := len(matrix[0])
	if cols == 0 {
		return newValidationError("Matrix must contain at least one row and one column")
	}

	for _, row := range matrix {
		if row == nil || len(row) == 0 {
			return newValidationError("Matrix must contain at least one row and one column")
		}

		if len(row) != cols {
			return newValidationError("Matrix must be rectangular")
		}

		for _, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return newValidationError("Matrix must contain only numeric values")
			}
		}

	}

	return nil
}
