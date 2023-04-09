package models

import "time"

type CryptoPrice struct {
	ID           int64
	CryptoName   string
	FiatName     string
	CurrentPrice float64
	MinPrice     float64
	MaxPrice     float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
