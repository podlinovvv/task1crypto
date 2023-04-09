package main

import (
	"log"
	"os"
	"task1crypto/internal/api"
	"task1crypto/internal/cache"
	"task1crypto/internal/db"
	"task1crypto/internal/telegram"
	"time"
)

const (
	updateInterval = 5 * time.Minute
)

func main() {
	apiKey := "226344ed8e181b034282e7ac12ca4170497b413abe9c3f8bccd5d7c7c6a535b0"
	botToken := "6200298459:AAFY4TZdDWshAgNwJCazc4gnclg7PwAgGqo"

	//apiKey := os.Getenv("apiKey")
	//botToken := os.Getenv("botToken")
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")

	dbClient, err := db.NewDBClient(dbConnStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	err = dbClient.CreateTable()
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	apiClient := api.NewCryptoCompareClient(apiKey)

	bot, err := telegram.NewBot(botToken)
	if err != nil {
		log.Fatalf("failed to create Telegram bot: %v", err)
	}

	pricesCache := cache.NewCache(dbClient)
	prices, _ := apiClient.GetCryptoRates()
	pricesCache.AddAndRecalculate(prices)

	ticker := time.NewTicker(updateInterval)
	go func() {
		for range ticker.C {
			prices, _ = apiClient.GetCryptoRates()
			pricesCache.AddAndRecalculate(prices)
		}
	}()

	bot.Start(pricesCache)
}

// мы передаём клиент бд в кэш, а кэш передаём в бота
// кажется, это как-то неправильно
