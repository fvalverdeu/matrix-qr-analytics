package services

import (
	"errors"
	"math"
	"testing"
)

func TestMatrixValidator_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		matrix      [][]float64
		wantErr     bool
		wantMessage string
	}{
		{
			name:        "nil_matrix",
			matrix:      nil,
			wantErr:     true,
			wantMessage: "Matrix is required",
		},
		{
			name:        "empty_matrix",
			matrix:      [][]float64{},
			wantErr:     true,
			wantMessage: "Matrix cannot be empty",
		},
		{
			name: "empty_first_row",
			matrix: [][]float64{
				{},
			},
			wantErr:     true,
			wantMessage: "Matrix must contain at least one row and one column",
		},
		{
			name: "empty_later_row",
			matrix: [][]float64{
				{1, 2},
				{},
			},
			wantErr:     true,
			wantMessage: "Matrix must contain at least one row and one column",
		},
		{
			name: "nil_later_row",
			matrix: [][]float64{
				{1, 2},
				nil,
			},
			wantErr:     true,
			wantMessage: "Matrix must contain at least one row and one column",
		},
		{
			name: "ragged_rows",
			matrix: [][]float64{
				{1, 2},
				{3},
			},
			wantErr:     true,
			wantMessage: "Matrix must be rectangular",
		},
		{
			name: "contains_nan",
			matrix: [][]float64{
				{1, math.NaN()},
			},
			wantErr:     true,
			wantMessage: "Matrix must contain only numeric values",
		},
		{
			name: "contains_positive_infinity",
			matrix: [][]float64{
				{1, math.Inf(1)},
			},
			wantErr:     true,
			wantMessage: "Matrix must contain only numeric values",
		},
		{
			name: "contains_negative_infinity",
			matrix: [][]float64{
				{1, math.Inf(-1)},
			},
			wantErr:     true,
			wantMessage: "Matrix must contain only numeric values",
		},
		{
			name: "valid_scalar_1x1",
			matrix: [][]float64{
				{5},
			},
			wantErr: false,
		},
		{
			name: "valid_square",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
			},
			wantErr: false,
		},
		{
			name: "valid_tall",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			wantErr: false,
		},
		{
			name: "valid_wide",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
			},
			wantErr: false,
		},
		{
			name: "valid_negative_decimal_values",
			matrix: [][]float64{
				{-1.5, 2.25},
				{3.75, -4.5},
			},
			wantErr: false,
		},
		{
			name: "valid_zero_matrix",
			matrix: [][]float64{
				{0, 0},
				{0, 0},
			},
			wantErr: false,
		},
	}

	validator := NewMatrixValidator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.matrix)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("expected ValidationError type, got %T", err)
				}

				if validationErr.Message != tc.wantMessage {
					t.Fatalf("unexpected error message: got %q, want %q", validationErr.Message, tc.wantMessage)
				}

				if err.Error() != tc.wantMessage {
					t.Fatalf("unexpected error string: got %q, want %q", err.Error(), tc.wantMessage)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
