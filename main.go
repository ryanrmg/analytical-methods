package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ryanrmg/analytical-methods/analytics"
	projectx "github.com/ryanrmg/projectx-api"
)

func main() {
	username := os.Getenv("PROJECTX_USERNAME")
	apiKey := os.Getenv("PROJECTX_API_KEY")

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		log.Fatal("Skipping integration test (no credentials provided)")
	}

	client := projectx.NewProjectXClient(
		"https://api.topstepx.com/api",
		"https://rtc.topstepx.com/hubs/",
		username,
		apiKey,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	analyticsService := analytics.NewAnalyticsService(client)

	analyticsService.LogMarket.StreamTradesToCSV(ctx)

}
