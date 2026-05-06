package analytics

import (
	"github.com/ryanrmg/projectx-api"
)

type AnalyticsService struct {
	LogUser   *LogServiceUser
	LogMarket *LogServiceMarket
	Pairs     *PairsIndicator
}

func NewAnalyticsService(client *projectx.ProjectXClient) *AnalyticsService {
	a := &AnalyticsService{}
	a.LogUser = &LogServiceUser{client: client}
	a.LogMarket = &LogServiceMarket{client: client}
	a.Pairs = &PairsIndicator{client: client, window: 20}
	return a
}
