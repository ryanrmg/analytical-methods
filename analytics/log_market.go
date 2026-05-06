package analytics

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ryanrmg/projectx-api"
)

type LogServiceMarket struct {
	client *projectx.ProjectXClient
}

func (l *LogServiceMarket) StreamQuotesToCSV(ctx context.Context) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("quotes-stream-market-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayQuoteCSVHeader() + "\n")

	for {
		select {
		case order := <-l.client.Realtime.QuotesStream():
			f.WriteString(order.ToCSVRow() + "\n")

		case <-ctx.Done():
			fmt.Println("stopped logging quotes")
			return nil
		}
	}
}

func (l *LogServiceMarket) StreamTradesToCSV(ctx context.Context) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("trade-stream-market-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayTradeCSVHeader() + "\n")

	for {
		select {
		case order := <-l.client.Realtime.TradesStream():
			f.WriteString(order.ToCSVRow() + "\n")

		case <-ctx.Done():
			fmt.Println("stopped logging trades")
			return nil
		}
	}
}

func (l *LogServiceMarket) StreamDepthToCSV(ctx context.Context) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("depth-stream-market-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayDepthCSVHeader() + "\n")

	for {
		select {
		case order := <-l.client.Realtime.DepthStream():
			f.WriteString(order.ToCSVRow() + "\n")

		case <-ctx.Done():
			fmt.Println("stopped logging depth")
			return nil
		}
	}
}
