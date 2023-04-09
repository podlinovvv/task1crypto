package cache

import (
	"log"
	"task1crypto/internal/db"
)

type Cache struct {
	DBClient *db.DBClient
	data     map[string]map[string][]float64
}

func NewCache(DBClient *db.DBClient) *Cache {
	return &Cache{
		DBClient: DBClient,
		data:     make(map[string]map[string][]float64),
	}
}

func (c *Cache) AddAndRecalculate(newDataframe map[string]map[string]float64) {
	pricesWereUpdated := false

	for cryptoName, valuesMap := range newDataframe {
		if _, ok := c.data[cryptoName]; !ok {
			c.data[cryptoName] = make(map[string][]float64)
		}

		for fiatName, price := range valuesMap {
			if _, ok := c.data[cryptoName][fiatName]; !ok {
				c.data[cryptoName][fiatName] = make([]float64, 3)
			}

			c.data[cryptoName][fiatName][1] = price
			min, max := c.data[cryptoName][fiatName][0], c.data[cryptoName][fiatName][2]

			if price > max {
				max = price
				pricesWereUpdated = true
			}
			if price < min || min == 0 {
				min = price
				pricesWereUpdated = true
			}
			c.data[cryptoName][fiatName][0] = min
			c.data[cryptoName][fiatName][2] = max

		}
	}

	if pricesWereUpdated {
		c.updateToDB()
	}

}

func (c *Cache) GetData() map[string]map[string][]float64 {
	if len(c.data) == 0 {
		c.updateFromDB()
	}
	return c.data
}

func (c *Cache) updateFromDB() {
	if len(c.data) == 0 {
		data, err := c.DBClient.GetCryptoPrices()
		if err != nil {
			log.Fatal(err)
		}
		c.data = data
	}
}

func (c *Cache) updateToDB() error {
	// обновить данные в бд
	err := c.DBClient.InsertCryptoPrices(c.data)
	if err != nil {
		log.Fatalf("failed to insert data into the database: %v", err)
		return err
	}
	return nil
}
