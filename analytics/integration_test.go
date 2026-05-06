package analytics

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ryanrmg/projectx-api"
)

var (
	username = os.Getenv("PROJECTX_USERNAME")
	apiKey   = os.Getenv("PROJECTX_API_KEY")
)

func TestAnalytics_Live(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := projectx.NewProjectXClient(
		"https://api.topstepx.com/api",
		"rtc.topstepx.com/hubs/",
		username,
		apiKey,
	)

	// connect
	if err := client.Realtime.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}

	// subscribe
	contract1 := "CON.F.US.MNQ.M26"
	symbolId1 := "F.US.MNQ"
	if err := client.Realtime.SubscribeContractTrades(contract1); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	contract2 := "CON.F.US.MES.M26"
	symbolId2 := "F.US.MES"
	if err := client.Realtime.SubscribeContractTrades(contract2); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	analyticsService := NewAnalyticsService(client)
	corrStream := analyticsService.Pairs.Correlation(ctx, symbolId1, symbolId2)

	timeout := time.After(25 * time.Second)
	received := false

	for {
		select {

		case c := <-corrStream:
			log.Println("CORR:", c)
			received = true

			// once we get a value, test succeeded
			return

		case <-timeout:
			if !received {
				t.Fatal("did not receive correlation value in time")
			}
			return

		case <-ctx.Done():
			t.Fatal("context cancelled before receiving correlation")
		}
	}

}
