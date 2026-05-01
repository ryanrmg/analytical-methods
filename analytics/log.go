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

func (l *LogService) GetAndLogOrders(ctx context.Context, accountId int, startTime, endTime string) {
	orders, err := l.client.Orders.OrderSearch(ctx, projectx.OrderSearchRequest{
		AccountId:      accountId,
		StartTimestamp: startTime,
		EndTimestamp:   endTime,
	})

	if err != nil {
		log.Fatal("Failed to get orders")
	}

	fmt.Println(orders)
}
