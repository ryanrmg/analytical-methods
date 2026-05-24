package analytics

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ryanrmg/projectx-api"
)

type LogServiceUser struct {
	client *projectx.ProjectXClient
}

func (l *LogServiceUser) StreamOrdersToCSV(ctx context.Context) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("orders-stream-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayUserOrderCSVHeader() + "\n")

	for {
		select {
		case order := <-l.client.Realtime.UserOrdersStream():
			f.WriteString(order.ToCSVRow() + "\n")

		case <-ctx.Done():
			fmt.Println("stopped logging orders")
			return nil
		}
	}
}

func (l *LogServiceUser) StreamTradesToCSV(ctx context.Context) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("trade-stream-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayUserTradeCSVHeader() + "\n")

	for {
		select {
		case order := <-l.client.Realtime.UserTradeStream():
			f.WriteString(order.ToCSVRow() + "\n")

		case <-ctx.Done():
			fmt.Println("stopped logging trades")
			return nil
		}
	}
}

func (l *LogServiceUser) StreamPositionsToCSV(ctx context.Context) error {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("position-stream-%s.csv", startTime)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(projectx.GatewayUserPositionCSVHeader() + "\n")

	for {
		select {
		case order := <-l.client.Realtime.UserPositionStream():
			f.WriteString(order.ToCSVRow() + "\n")

		case <-ctx.Done():
			fmt.Println("stopped logging position")
			return nil
		}
	}
}
