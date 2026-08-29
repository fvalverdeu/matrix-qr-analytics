package services

import (
	"errors"
	"reflect"
	"testing"

	"qr-service-go/internal/models"
)

type fakeStatisticsClient struct {
	callCount int
	lastQ     [][]float64
	lastR     [][]float64
	result    models.Statistics
	err       error
}

func (f *fakeStatisticsClient) CalculateStatistics(q, r [][]float64) (models.Statistics, error) {
	f.callCount++
	f.lastQ = q
	f.lastR = r

	if f.err != nil {
		return models.Statistics{}, f.err
	}

	return f.result, nil
}

func TestQRServiceProcess_Success(t *testing.T) {
	statisticsResult := models.Statistics{
		Max:               9.25,
		Min:               -3.5,
		Average:           1.125,
		Sum:               6.75,
		HasDiagonalMatrix: false,
	}

	fakeClient := &fakeStatisticsClient{result: statisticsResult}
	service := NewQRService(NewMatrixValidator(), fakeClient)

	input := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}

	response, err := service.Process(input)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertMatrixShape(t, "Q", response.Q, 2, 2)
	assertMatrixShape(t, "R", response.R, 2, 3)

	if response.Statistics != statisticsResult {
		t.Fatalf("unexpected statistics: got %+v, want %+v", response.Statistics, statisticsResult)
	}

	if fakeClient.callCount != 1 {
		t.Fatalf("expected statistics client to be called once, got %d", fakeClient.callCount)
	}

	if !reflect.DeepEqual(fakeClient.lastQ, response.Q) {
		t.Fatalf("statistics client Q input differs from response Q")
	}

	if !reflect.DeepEqual(fakeClient.lastR, response.R) {
		t.Fatalf("statistics client R input differs from response R")
	}
}

func TestQRServiceProcess_InvalidMatrixDoesNotCallStatistics(t *testing.T) {
	fakeClient := &fakeStatisticsClient{}
	service := NewQRService(NewMatrixValidator(), fakeClient)

	response, err := service.Process(nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError type, got %T", err)
	}

	if fakeClient.callCount != 0 {
		t.Fatalf("expected statistics client call count 0, got %d", fakeClient.callCount)
	}

	if !reflect.DeepEqual(response, models.QRResponse{}) {
		t.Fatalf("expected zero-value QRResponse, got %+v", response)
	}
}

func TestQRServiceProcess_StatisticsErrorIsPropagated(t *testing.T) {
	sentinelErr := errors.New("statistics downstream failure")
	fakeClient := &fakeStatisticsClient{err: sentinelErr}
	service := NewQRService(NewMatrixValidator(), fakeClient)

	input := [][]float64{
		{1, 2},
		{3, 4},
	}

	response, err := service.Process(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, sentinelErr) {
		t.Fatalf("expected propagated sentinel error, got %v", err)
	}

	if fakeClient.callCount != 1 {
		t.Fatalf("expected statistics client to be called once, got %d", fakeClient.callCount)
	}

	if len(fakeClient.lastQ) == 0 || len(fakeClient.lastQ[0]) == 0 {
		t.Fatalf("expected non-empty Q to be passed to statistics client")
	}

	if len(fakeClient.lastR) == 0 || len(fakeClient.lastR[0]) == 0 {
		t.Fatalf("expected non-empty R to be passed to statistics client")
	}

	if !reflect.DeepEqual(response, models.QRResponse{}) {
		t.Fatalf("expected zero-value QRResponse on statistics error, got %+v", response)
	}
}

func assertMatrixShape(t *testing.T, name string, matrix [][]float64, expectedRows, expectedCols int) {
	t.Helper()

	if len(matrix) != expectedRows {
		t.Fatalf("%s row count mismatch: got %d, want %d", name, len(matrix), expectedRows)
	}

	for i := range matrix {
		if len(matrix[i]) != expectedCols {
			t.Fatalf("%s col count mismatch at row %d: got %d, want %d", name, i, len(matrix[i]), expectedCols)
		}
	}
}
