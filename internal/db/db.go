package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"task1crypto/internal/models"
)

type DBClient struct {
	db *sqlx.DB
}

func NewDBClient(connStr string) (*DBClient, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	return &DBClient{db: db}, nil
}

func (d *DBClient) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS crypto_prices (
		id SERIAL PRIMARY KEY,
		crypto_name VARCHAR(10) NOT NULL,
		fiat_name VARCHAR(10) NOT NULL,
		current_price NUMERIC(20, 8) NOT NULL,
		min_price NUMERIC(20, 8) NOT NULL,
		max_price NUMERIC(20, 8) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	`
	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func (d *DBClient) InsertCryptoPrices(prices map[string]map[string][]float64) error {
	query := `
	INSERT INTO crypto_prices (crypto_name, fiat_name, current_price, min_price, max_price, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	now := time.Now()

	for cryptoName, valuesMap := range prices {
		for fiatName, priceData := range valuesMap {
			_, err := d.db.Exec(query, cryptoName, fiatName, priceData[1], priceData[0], priceData[2], now, now)
			if err != nil {
				return fmt.Errorf("failed to insert data into crypto_prices: %v", err)
			}
		}
	}

	return nil
}

func (d *DBClient) GetCryptoPrices() (map[string]map[string][]float64, error) {
	query := `
	SELECT crypto_name, fiat_name, current_price, min_price, max_price
	FROM crypto_prices;
	`
	rows, err := d.db.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query crypto_prices: %v", err)
	}
	defer rows.Close()

	prices := make(map[string]map[string][]float64)
	for rows.Next() {
		var cp models.CryptoPrice
		err := rows.StructScan(&cp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if _, ok := prices[cp.CryptoName]; !ok {
			prices[cp.CryptoName] = make(map[string][]float64)
		}
		prices[cp.CryptoName][cp.FiatName] = []float64{cp.MinPrice, cp.CurrentPrice, cp.MaxPrice}
	}

	return prices, nil
}
