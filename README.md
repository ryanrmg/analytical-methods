# analytical-methods

# Log 

# Gemini

# Database

```go
func main() {
	ctx := context.Background()
	connStr := "postgres://trading_user:your_secure_password@localhost:5432/trading_db"

	// Connect to cluster
	pool, _ := pgxpool.New(ctx, connStr)
	defer pool.Close()

	// Initialize store layer
	store := NewDBStore(pool)

	// 1. Create the table
	if err := store.CreateResponseCacheTable(ctx); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// 2. Insert mock data (Simulating an API hit)
	endpoint := "/v1/market/eth"
	mockData := `{"symbol": "ETHUSDT", "price": "3450.10"}`
	
	if err := store.InsertResponse(ctx, endpoint, mockData); err != nil {
		log.Printf("Error inserting: %v", err)
	}

	// 3. Read data back out
	cachedJSON, fetchedAt, err := store.GetLatestResponse(ctx, endpoint)
	if err == nil {
		fmt.Printf("Found cached data from %s: %s\n", fetchedAt.Format(time.Kitchen), cachedJSON)
	}
}
```


```bash
psql -U trading_user -d trading_db -h localhost
```
