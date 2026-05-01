package analytics

import (
	"context"
	"fmt"

	"github.com/ryanrmg/projectx-api"
)

type LogService struct {
	client *projectx.ProjectXClient
}


func (l *LogService) GetAndLogOrders(ctx context.Context, accountId int, startTime, endTime string) {
	orders := l.client.Orders.OrderSearch(ctx, projectx.OrderSearchRequest{
		AccountId: accountId,
		StartTime: startTime,
		EndTime: endTime
	})

	fmt.Println(orders)
}