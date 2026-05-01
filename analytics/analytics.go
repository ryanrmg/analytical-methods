package analytics

import (
	"github.com/ryanrmg/projectx-api"
)

type AnalyticsService struct {
	Log *LogService
}

func NewAnalyticsService(client *projectx.ProjectXClient) *AnalyticsService {
	a := &AnalyticsService{}
	a.Log = &LogService{client: client}
	return a
}
