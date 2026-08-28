package services

import (
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

	var factorization mat.QR
	factorization.Factorize(a)

	var q mat.Dense
	var r mat.Dense
	factorization.QTo(&q)
	factorization.RTo(&r)

	return denseToSlice(&q), denseToSlice(&r), nil
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
