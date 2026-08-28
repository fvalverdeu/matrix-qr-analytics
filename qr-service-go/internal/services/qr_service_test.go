package services

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

const tolerance = 1e-9

func TestQRServiceDecompose_CharacterizationInvariants(t *testing.T) {
	testCases := []struct {
		name   string
		matrix [][]float64
	}{
		{
			name: "scalar_1x1",
			matrix: [][]float64{
				{5},
			},
		},
		{
			name: "square_2x2",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
			},
		},
		{
			name: "tall_3x2",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			name: "rank_deficient",
			matrix: [][]float64{
				{1, 2},
				{2, 4},
				{3, 6},
			},
		},
		{
			name: "zero_matrix",
			matrix: [][]float64{
				{0, 0},
				{0, 0},
				{0, 0},
			},
		},
		{
			name: "negative_decimal_values",
			matrix: [][]float64{
				{-1.5, 2.25},
				{3.75, -4.5},
				{-6.125, 7.0},
			},
		},
	}

	service := NewQRService(NewMatrixValidator(), nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qSlice, rSlice, err := service.decompose(tc.matrix)
			if err != nil {
				t.Fatalf("decompose returned error: %v", err)
			}

			m := len(tc.matrix)
			n := len(tc.matrix[0])

			assertSliceDimensions(t, "Q", qSlice, m, m)
			assertSliceDimensions(t, "R", rSlice, m, n)

			qDense := denseFromSlice(qSlice)
			rDense := denseFromSlice(rSlice)
			aDense := denseFromSlice(tc.matrix)

			assertOrthogonality(t, qDense, m)
			assertReconstruction(t, qDense, rDense, aDense)
			assertUpperTrapezoidal(t, rDense)
			assertAllFinite(t, "Q", qDense)
			assertAllFinite(t, "R", rDense)
		})
	}
}

func denseFromSlice(values [][]float64) *mat.Dense {
	rows := len(values)
	cols := len(values[0])
	data := make([]float64, 0, rows*cols)
	for _, row := range values {
		data = append(data, row...)
	}

	return mat.NewDense(rows, cols, data)
}

func assertSliceDimensions(t *testing.T, matrixName string, values [][]float64, expectedRows, expectedCols int) {
	t.Helper()

	if len(values) != expectedRows {
		t.Fatalf("%s row count mismatch: got %d, want %d", matrixName, len(values), expectedRows)
	}

	for i := range values {
		if len(values[i]) != expectedCols {
			t.Fatalf("%s col count mismatch at row %d: got %d, want %d", matrixName, i, len(values[i]), expectedCols)
		}
	}
}

func assertOrthogonality(t *testing.T, q *mat.Dense, m int) {
	t.Helper()

	var qtq mat.Dense
	qtq.Mul(q.T(), q)

	identity := mat.NewDense(m, m, nil)
	for i := 0; i < m; i++ {
		identity.Set(i, i, 1)
	}

	var diff mat.Dense
	diff.Sub(&qtq, identity)

	err := frobeniusNorm(&diff)
	threshold := tolerance * math.Max(1, float64(m))
	if err > threshold {
		t.Fatalf("orthogonality invariant failed: ||Q^TQ - I||_F=%g, threshold=%g", err, threshold)
	}
}

func assertReconstruction(t *testing.T, q, r, a *mat.Dense) {
	t.Helper()

	var reconstructed mat.Dense
	reconstructed.Mul(q, r)

	var diff mat.Dense
	diff.Sub(&reconstructed, a)

	err := frobeniusNorm(&diff)
	aNorm := frobeniusNorm(a)
	threshold := tolerance * math.Max(1, aNorm)
	if err > threshold {
		t.Fatalf("reconstruction invariant failed: ||Q*R - A||_F=%g, threshold=%g", err, threshold)
	}
}

func assertUpperTrapezoidal(t *testing.T, r *mat.Dense) {
	t.Helper()

	rows, cols := r.Dims()
	for i := 0; i < rows; i++ {
		limit := i
		if limit > cols {
			limit = cols
		}

		for j := 0; j < limit; j++ {
			if math.Abs(r.At(i, j)) > tolerance {
				t.Fatalf("upper trapezoidal invariant failed at R[%d,%d]=%g", i, j, r.At(i, j))
			}
		}
	}
}

func assertAllFinite(t *testing.T, matrixName string, matrix *mat.Dense) {
	t.Helper()

	rows, cols := matrix.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			value := matrix.At(i, j)
			if math.IsNaN(value) || math.IsInf(value, 0) {
				t.Fatalf("%s contains non-finite value at [%d,%d]=%v", matrixName, i, j, value)
			}
		}
	}
}

func frobeniusNorm(matrix mat.Matrix) float64 {
	rows, cols := matrix.Dims()
	var sum float64
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			value := matrix.At(i, j)
			sum += value * value
		}
	}

	return math.Sqrt(sum)
}
