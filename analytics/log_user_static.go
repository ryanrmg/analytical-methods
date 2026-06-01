package analytics

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ryanrmg/projectx-api"
)

type StaticLogServiceUser struct {
	client *projectx.ProjectXClient
	db     *DBStore
}

func (l *StaticLogServiceUser) LogTradesToCSV(ctx context.Context, accountId int) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("trade-stream-market-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayTradeCSVHeader() + "\n")

	trades, err := l.client.Trades.Search(ctx, projectx.TradeSearchRequest{
		AccountId:      accountId,
		StartTimestamp: time.Now().UTC().Add(-time.Duration(24*7) * time.Hour).Format(time.RFC3339),
		EndTimestamp:   time.Now().UTC().Format(time.RFC3339),
	})

	if err != nil {
		return err
	}

	for _, trade := range trades {
		f.WriteString(trade.ToCSVRow() + "\n")
	}
	return nil
}

func (l *StaticLogServiceUser) LogTradesToDB(ctx context.Context, accountId int) error {
	trades, err := l.client.Trades.Search(ctx, projectx.TradeSearchRequest{
		AccountId:      accountId,
		StartTimestamp: time.Now().UTC().Add(-time.Duration(24*7) * time.Hour).Format(time.RFC3339),
		EndTimestamp:   time.Now().UTC().Format(time.RFC3339),
	})

	if err != nil {
		return err
	}

	for _, trade := range trades {
		err := l.db.SaveUserTrade(context.Background(), trade)
		if err != nil {
			fmt.Printf("Error caching trade %d to database: %v", trade.Id, err)
		}
	}
	return nil
}
