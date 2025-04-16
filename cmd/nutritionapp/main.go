package main

import (
	"fmt"
	"log"
	"nutritionapp/pkg/client"
	"nutritionapp/pkg/db"
	"nutritionapp/pkg/fdc"
	"nutritionapp/pkg/server"
	"os"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("FDC_API_KEY")
	if apiKey == "" {
		log.Fatal("FDC_API_KEY environment variable not set")
	}

	// Initialize SQLite database
	sqliteDB, err := db.NewSQLiteDB("nutritionapp.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create request channel
	requests := make(chan server.Request)

	// Start server
	srv := server.NewServer(sqliteDB, fdc.NewFoodProcessor(apiKey))
	go srv.Start()

	// Start client
	cli := client.NewClient(requests)
	fmt.Println("Starting NutritionApp...")
	cli.Start()
}
