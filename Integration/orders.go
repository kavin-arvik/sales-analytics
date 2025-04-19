package integration

import (
	db "SalesAnalytics/DB"
	"SalesAnalytics/constants"
	"SalesAnalytics/handlers"
	"SalesAnalytics/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func CustomerAnalysis(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("r.Method  - ", r.Method)

	log.Println("CustomerAnalysis + ")

	var lResp models.RESPONSE
	if r.Method == "GET" {

		lResp.Status = "S"
		lResp.Msg = "Data fetched successfully!!"
		lVars := mux.Vars(r)
		lId := lVars["id"]
		var lErr error

		if lId == "dateRange" {
			lResp.Customers, lErr = TotCustomers(r)
		} else if lId == "orders" {
			lResp.Customers, lErr = TotOrders(r)
		} else if lId == "avgOrderValue" {
			lResp.Customers, lErr = TotAvgOrdValue(r)
		}

		//Handle overall errors
		if lErr != nil {
			lResp.Status = "E"
			lResp.Msg = "Error while fetching"
		}

	} else if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		lResp.Status = constants.ERROR
		lResp.Msg = "Invalid HTTP method"
	}
	lData, lErr := json.Marshal(lResp)
	if lErr != nil {
		fmt.Fprintf(w, "Error taking data"+lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}
	log.Println("CustomerAnalysis - ")
}

func TotCustomers(r *http.Request) (*models.CUSTOMERS, error) {
	log.Println("TopSaleDateRange + ")
	var lTotCustomers models.CUSTOMERS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				select count(*) 
				from orders o, order_items oi 
				where o.id = oi.order_id
				and o.date_of_sale BETWEEN $1 AND $2`

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TC01) ", lErr)
		return &lTotCustomers, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lTotCustomers.TotCustomers)
			if lErr != nil {
				log.Println("Error (TC02) ", lErr)
				return &lTotCustomers, lErr
			}
		}
	}
	log.Println("TopSaleDateRange - ")
	return &lTotCustomers, nil
}

func TotOrders(r *http.Request) (*models.CUSTOMERS, error) {
	log.Println("TotOrders + ")
	var lTotCustomers models.CUSTOMERS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				select count(*) from orders o, order_items oi 
				where o.id = oi.order_id
				and o.date_of_sale BETWEEN $1 AND $2`

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TO01) ", lErr)
		return &lTotCustomers, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lTotCustomers.TotOrders)
			if lErr != nil {
				log.Println("Error (TO02) ", lErr)
				return &lTotCustomers, lErr
			}
		}
	}
	log.Println("TotOrders - ")
	return &lTotCustomers, nil
}

func TotAvgOrdValue(r *http.Request) (*models.CUSTOMERS, error) {
	log.Println("TotAvgOrdValue + ")
	var lTotCustomers models.CUSTOMERS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				SELECT ROUND(AVG(order_total), 2) AS average_order_value
				FROM (
					SELECT o.id, 
						SUM((oi.quantity * oi.unit_price) + oi.discount + oi.shipping_cost) AS order_total
					FROM orders o
					JOIN order_items oi ON o.id = oi.order_id
				WHERE o.date_of_sale BETWEEN $1 AND $2
					GROUP BY o.id
				) AS order_totals;`

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TAOV01) ", lErr)
		return &lTotCustomers, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lTotCustomers.AvgOrderValue)
			if lErr != nil {
				log.Println("Error (TAOV02) ", lErr)
				return &lTotCustomers, lErr
			}
		}
	}
	log.Println("TotAvgOrdValue - ")
	return &lTotCustomers, nil
}
