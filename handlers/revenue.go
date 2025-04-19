package handlers

import (
	db "SalesAnalytics/DB"
	"SalesAnalytics/constants"
	"SalesAnalytics/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RevenueHandler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("r.Method  - ", r.Method)

	log.Println("RevenueHandler + ")

	var lResp models.RESPONSE
	if r.Method == "GET" {

		lResp.Status = "S"
		lResp.Msg = "Data fetched successfully!!"
		lVars := mux.Vars(r)
		lId := lVars["id"]
		var lErr error

		if lId == "dateRange" {
			lResp.TotRevenue, lErr = RevDateRange(r)
		} else if lId == "product" {
			lResp.TotProdRevenue, lErr = RevProducts(r)
		} else if lId == "category" {
			lResp.TotCategRevenue, lErr = RevCategory(r)
		} else if lId == "region" {
			lResp.TotRegionRevenue, lErr = RevRegion(r)
		}

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
	// json.NewEncoder(w).Encode(map[string]float64{"total_revenue": lRevenue})
	log.Println("RevenueHandler - ")
}

func GetDateRange(r *http.Request) (time.Time, time.Time) {
	log.Println("getDateRange + ")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)
	log.Println("getDateRange - ")
	return startDate, endDate
}

func RevDateRange(r *http.Request) (*float64, error) {
	log.Println("RevDateRange + ")
	startDate, endDate := GetDateRange(r)

	var lRevenue float64
	lSelectString := `
        SELECT SUM((oi.unit_price * oi.quantity))
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        WHERE o.date_of_sale BETWEEN $1 AND $2
    `
	lErr := db.GDBConnection.QueryRow(lSelectString, startDate, endDate).Scan(&lRevenue)
	if lErr != nil {
		log.Println("Error (HRRDG01) ", lErr)
		return &lRevenue, lErr
	}
	log.Println("RevDateRange - ")
	return &lRevenue, nil
}

func RevProducts(r *http.Request) (*[]models.REVPRODUCTS, error) {
	log.Println("RevProducts + ")
	var lRevProductsRec models.REVPRODUCTS
	var lRevProductsArr []models.REVPRODUCTS

	startDate, endDate := GetDateRange(r)

	lSelectString := `
        select name, sum((quantity * unit_price)) 
		from orders o, order_items oi , products p 
		where o.id = oi.order_id and oi.product_id = p.id  
			and o.date_of_sale BETWEEN $1 AND $2
		group by name
    `
	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (HRRP01) ", lErr)
		return &lRevProductsArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lRevProductsRec.Name, &lRevProductsRec.Revenue)
			if lErr != nil {
				log.Println("Error (HRRP02) ", lErr)
				return &lRevProductsArr, lErr
			} else {
				lRevProductsArr = append(lRevProductsArr, lRevProductsRec)
			}
		}
	}
	log.Println("RevProducts - ")
	return &lRevProductsArr, nil
}

func RevCategory(r *http.Request) (*[]models.REVCATEGORY, error) {
	log.Println("RevCategory + ")
	startDate, endDate := GetDateRange(r)
	var lRevCategoryRec models.REVCATEGORY
	var lRevCategoryArr []models.REVCATEGORY

	lSelectString := `
        select category, sum((oi.quantity * oi.unit_price)) 
		from orders o, order_items oi , products p 
		where o.id = oi.order_id and oi.product_id = p.id  
			and o.date_of_sale BETWEEN $1 AND $2
		group by p.category 
    `

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (HRRC01) ", lErr)
		return &lRevCategoryArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lRevCategoryRec.Category, &lRevCategoryRec.Revenue)
			if lErr != nil {
				log.Println("Error (HRRC02) ", lErr)
				return &lRevCategoryArr, lErr
			} else {
				log.Println("lRevCategoryRec - ", lRevCategoryRec)
				lRevCategoryArr = append(lRevCategoryArr, lRevCategoryRec)
			}
		}
	}
	log.Println("RevCategory - ")
	return &lRevCategoryArr, nil
}

func RevRegion(r *http.Request) (*[]models.REVREGION, error) {
	log.Println("RevRegion + ")
	startDate, endDate := GetDateRange(r)
	var lRevRegionRec models.REVREGION
	var lRevRegionArr []models.REVREGION

	lSelectString := `
			    select region, sum((oi.quantity * oi.unit_price)) 
				from orders o, order_items oi
				where o.id = oi.order_id 
					and o.date_of_sale BETWEEN $1 AND $2
				group by o.region 
    `
	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (HRR01) ", lErr)
		return &lRevRegionArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lRevRegionRec.Region, &lRevRegionRec.Revenue)
			if lErr != nil {
				log.Println("Error (HRR02) ", lErr)
				return &lRevRegionArr, lErr
			} else {
				log.Println("lRevCategoryRec - ", lRevRegionRec)
				lRevRegionArr = append(lRevRegionArr, lRevRegionRec)
			}
		}
	}
	log.Println("RevRegion - ")
	return &lRevRegionArr, nil
}
