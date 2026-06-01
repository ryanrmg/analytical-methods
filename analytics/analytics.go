package analytics

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryanrmg/projectx-api"
)

type AnalyticsService struct {
	LogUser       *LogServiceUser
	LogUserStatic *StaticLogServiceUser
	LogMarket     *LogServiceMarket
	Pairs         *PairsIndicator
	DataBase      *DBStore
}

// Pass the *pgxpool.Pool here so the service layer can reach the database
func NewAnalyticsService(client *projectx.ProjectXClient) *AnalyticsService {

	useDatabase := os.Getenv("USE_DB") == "true"

	if useDatabase {
		username := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")

		connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/trading_db", username, password)

		ctx := context.Background()
		// Connect to cluster
		dbPool, _ := pgxpool.New(ctx, connStr)
		defer dbPool.Close()

		// Initialize your DBStore wrapper with the pool
		dbStore := NewDBStore(dbPool)

		if err := dbStore.CreateUserTradeTable(ctx); err != nil {
			log.Printf("Warning: failed to verify/create trades table: %v", err)
		}

		return &AnalyticsService{
			LogUser:       &LogServiceUser{client: client},
			LogUserStatic: &StaticLogServiceUser{client: client, db: dbStore},
			LogMarket:     &LogServiceMarket{client: client},
			Pairs:         &PairsIndicator{client: client, window: 120},
			DataBase:      dbStore,
		}
	}

	return &AnalyticsService{
		LogUser:       &LogServiceUser{client: client},
		LogUserStatic: &StaticLogServiceUser{client: client, db: nil},
		LogMarket:     &LogServiceMarket{client: client},
		Pairs:         &PairsIndicator{client: client, window: 120},
		DataBase:      nil,
	}
}
