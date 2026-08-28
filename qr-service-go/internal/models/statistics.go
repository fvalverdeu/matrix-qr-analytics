package models

type StatisticsRequest struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

type Statistics struct {
	Max               float64 `json:"max"`
	Min               float64 `json:"min"`
	Average           float64 `json:"average"`
	Sum               float64 `json:"sum"`
	HasDiagonalMatrix bool    `json:"hasDiagonalMatrix"`
}
