package analytics

import (
	"context"
	"log"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ryanrmg/projectx-api"
)

var (
	client *projectx.ProjectXClient
	once   sync.Once
)

func getClient() *projectx.ProjectXClient {
	once.Do(func() {
		client = projectx.NewProjectXClient(
			"https://api.topstepx.com/api",
			"rtc.topstepx.com/hubs/",
			os.Getenv("PROJECTX_USERNAME"),
			os.Getenv("PROJECTX_API_KEY"),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := client.Realtime.Connect(ctx); err != nil {
			panic(err)
		}
	})
	return client
}

func TestAnalytics_Static(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := getClient()

	// connect
	if err := client.Realtime.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}

	req1 := projectx.BarHistoryRequest{
		ContractId:        "CON.F.US.MNQ.M26",
		Live:              false,
		EndTime:           time.Now().UTC().Format(time.RFC3339),
		StartTime:         time.Now().UTC().Add(-time.Duration(24*5) * time.Hour).Format(time.RFC3339),
		Unit:              2,
		UnitNumber:        5,
		Limit:             100,
		IncludePartialBar: false,
	}

	req2 := projectx.BarHistoryRequest{
		ContractId:        "CON.F.US.MES.M26",
		Live:              false,
		EndTime:           time.Now().UTC().Format(time.RFC3339),
		StartTime:         time.Now().UTC().Add(-time.Duration(24*5) * time.Hour).Format(time.RFC3339),
		Unit:              2,
		UnitNumber:        5,
		Limit:             100,
		IncludePartialBar: false,
	}

	analyticsService := NewAnalyticsService(client)

	corr := analyticsService.Pairs.StaticCorrelation(ctx, req1, req2)

	log.Println("STATIC CORRELATION:", corr)

	// basic sanity assertions
	if math.IsNaN(corr) {
		t.Fatal("correlation is NaN")
	}

	if corr < -1 || corr > 1 {
		t.Fatalf("invalid correlation value: %f", corr)
	}

	// MNQ/MES should generally have positive correlation
	if corr < 0.3 {
		t.Logf("warning: unexpectedly low correlation: %f", corr)
	}
}

func TestAnalytics_Live(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := getClient()

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
