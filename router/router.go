package router

import (
	integration "SalesAnalytics/Integration"
	"SalesAnalytics/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/refresh", handlers.RefreshHandler).Methods("POST")
	r.HandleFunc("/api/revenue/{id}", handlers.RevenueHandler)
	r.HandleFunc("/api/nProducts/{id}", integration.TopNProducts)
	r.HandleFunc("/api/customers/{id}", integration.CustomerAnalysis)
	r.HandleFunc("/api/calculations/{id}", integration.SalesCalculations)

}
