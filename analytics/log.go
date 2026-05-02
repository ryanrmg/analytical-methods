package analytics

import (
	"context"
	"fmt"
	"log"

	"github.com/ryanrmg/projectx-api"
)

type LogService struct {
	client *projectx.ProjectXClient
}

func (l *LogService) StreamOrdersToCSV(ctx context.Context) error {
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
