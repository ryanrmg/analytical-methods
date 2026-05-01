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

	accounts, error := client.Accounts.Search(ctx, projectx.AccountSearchRequest{
		OnlyActiveAccounts: true,
	})
	if error != nil {
		log.Fatal("Failed to get accounts")
	}

	analyticsService := analytics.NewAnalyticsService(client)

	for _, account := range accounts {
		analyticsService.Log.GetAndLogOrders(
			ctx,
			account.Id, time.Now().UTC().Add(-time.Duration(24*5)*time.Hour).Format(time.RFC3339),
			time.Now().UTC().Format(time.RFC3339),
		)
	}

}
