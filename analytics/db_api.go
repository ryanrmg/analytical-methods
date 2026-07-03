package analytics

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// TradeHandler handles requests from the React frontend
type TradeHandler struct {
	DB *DBStore // Using the DBStore wrapper we created earlier
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	
	"github.com/jackc/pgx/v5"
)

func (h *TradeHandler) GetUserTradesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	accountIDStr := r.URL.Query().Get("accountId")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid accountId"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Try to read from your local PostgreSQL database first
	trades, err := h.DB.GetTradesByAccount(ctx, accountID)
	
	// 2. If nothing is available (empty slice) OR the database throws an error, fallback to API
	if err != nil || len(trades) == 0 {
		// Log that we missed the cache
		println("Database cache miss. Fetching fresh data from external Trading API...")

		// 3. Call your external API wrapper (e.g., via your analytics user service)
		// This is where you connect to the live projectx client
		apiTrades, fetchErr := h.Analytics.StaticLogServiceUser.LogTradesToDB(ctx, accountID)
		if fetchErr != nil {
			http.Error(w, `{"error": "Database empty and external API failed"}`, http.StatusInternalServerError)
			return
		}

		// 4. Populate the database asynchronously so it doesn't block the user's view
		go func(tradesToSave []GatewayUserTrade) {
			for _, trade := range tradesToSave {
				if saveErr := h.DB.SaveUserTrade(context.Background(), trade); saveErr != nil {
					// Simply log the error; don't crash the application
					println("Error populating database:", saveErr.Error()) 
				}
			}
			println("Successfully populated database with fresh API logs.")
		}(apiTrades)

		// 5. Assign the fresh data to return to the frontend immediately
		trades = apiTrades
	} else {
		println("Database cache hit! Returning stored logs instantly.")
	}

	// 6. Return the data (whether it came from DB or fresh API) to React
	json.NewEncoder(w).Encode(trades)
}
