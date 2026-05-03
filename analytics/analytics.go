package analytics

import (
	"github.com/ryanrmg/projectx-api"
)

type AnalyticsService struct {
	LogUser   *LogServiceUser
	LogMarket *LogServiceMarket
}

func NewAnalyticsService(client *projectx.ProjectXClient) *AnalyticsService {
	a := &AnalyticsService{}
	a.LogUser = &LogServiceUser{client: client}
	a.LogMarket = &LogServiceMarket{client: client}
	return a
}
