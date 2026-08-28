package services

import (
	"gonum.org/v1/gonum/lapack/lapack64"
	"gonum.org/v1/gonum/mat"

	"qr-service-go/internal/models"
)

type StatisticsClient interface {
	CalculateStatistics(q, r [][]float64) (models.Statistics, error)
}

type QRService struct {
	validator        *MatrixValidator
	statisticsClient StatisticsClient
}

func NewQRService(validator *MatrixValidator, statisticsClient StatisticsClient) *QRService {
	return &QRService{
		validator:        validator,
		statisticsClient: statisticsClient,
	}
}

func (s *QRService) Process(matrix [][]float64) (models.QRResponse, error) {
	q, r, err := s.decompose(matrix)
	if err != nil {
		return models.QRResponse{}, err
	}

	statistics, err := s.statisticsClient.CalculateStatistics(q, r)
	if err != nil {
		return models.QRResponse{}, err
	}

	return models.QRResponse{
		Q:          q,
		R:          r,
		Statistics: statistics,
	}, nil
}

func (s *QRService) decompose(matrix [][]float64) ([][]float64, [][]float64, error) {
	if err := s.validator.Validate(matrix); err != nil {
		return nil, nil, err
	}

	rows := len(matrix)
	cols := len(matrix[0])

	data := make([]float64, 0, rows*cols)
	for _, row := range matrix {
		data = append(data, row...)
	}

	a := mat.NewDense(rows, cols, data)
	factor := mat.DenseCopyOf(a)

	k := minInt(rows, cols)
	tau := make([]float64, k)

	work := []float64{0}
	lapack64.Geqrf(factor.RawMatrix(), tau, work, -1)
	work = make([]float64, int(work[0]))
	lapack64.Geqrf(factor.RawMatrix(), tau, work, len(work))

	r := extractR(factor, rows, cols)
	q := generateFullQ(factor, rows, k, tau)

	return denseToSlice(q), denseToSlice(r), nil
}

func extractR(factor *mat.Dense, rows, cols int) *mat.Dense {
	r := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := i; j < cols; j++ {
			r.Set(i, j, factor.At(i, j))
		}
	}

	return r
}

func generateFullQ(factor *mat.Dense, rows, k int, tau []float64) *mat.Dense {
	q := mat.NewDense(rows, rows, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < k; j++ {
			q.Set(i, j, factor.At(i, j))
		}
	}

	work := []float64{0}
	lapack64.Orgqr(q.RawMatrix(), tau, work, -1)
	work = make([]float64, int(work[0]))
	lapack64.Orgqr(q.RawMatrix(), tau, work, len(work))

	return q
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func denseToSlice(d *mat.Dense) [][]float64 {
	rows, cols := d.Dims()
	result := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = d.At(i, j)
		}
	}

	return result
}
