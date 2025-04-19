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

func TopNProducts(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "NVALUE, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("r.Method  - ", r.Method)

	log.Println("TopNProducts + ")

	var lResp models.RESPONSE
	if r.Method == "GET" {

		lNValue := r.Header.Get("NVALUE")

		lResp.Status = "S"
		lResp.Msg = "Data fetched successfully!!"
		lVars := mux.Vars(r)
		lId := lVars["id"]
		var lErr error

		if lId == "overall" {
			lResp.TopNProducts, lErr = TopSaleDateRange(r, lNValue)
		} else if lId == "category" {
			lResp.TopNProducts, lErr = TopSaleCategory(r, lNValue)
		} else if lId == "region" {
			lResp.TopNProducts, lErr = TopSaleRegion(r, lNValue)
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
	log.Println("TopNProducts - ")
}

func TopSaleDateRange(r *http.Request, pNValue string) (*[]models.TOPNPRODUCTS, error) {
	log.Println("TopSaleDateRange + ")
	var lTopNProdRec models.TOPNPRODUCTS
	var lTopNProdArr []models.TOPNPRODUCTS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				select *
			from (select p."name", sum(quantity) quantity
					from order_items oi, products p, orders o 
					where oi.product_id = p.id and o.id = oi.order_id
					and o.date_of_sale BETWEEN $1 AND $2
					group by p."name" ) t
					order by t.quantity desc
					limit ` + pNValue

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TSDR01) ", lErr)
		return &lTopNProdArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lTopNProdRec.Name, &lTopNProdRec.Qty)
			if lErr != nil {
				log.Println("Error (TSDR02) ", lErr)
				return &lTopNProdArr, lErr
			} else {
				lTopNProdArr = append(lTopNProdArr, lTopNProdRec)
			}
		}
	}
	log.Println("TopSaleDateRange - ")
	return &lTopNProdArr, nil
}

func TopSaleCategory(r *http.Request, pNValue string) (*[]models.TOPNPRODUCTS, error) {
	log.Println("TopSaleCategory + ")
	var lTopNProdRec models.TOPNPRODUCTS
	var lTopNProdArr []models.TOPNPRODUCTS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
					select *
				from (select p.category, sum(quantity) quantity
						from order_items oi, products p, orders o 
						where oi.product_id = p.id and o.id = oi.order_id
						and o.date_of_sale BETWEEN $1 AND $2
						group by p.category ) t
						order by t.quantity desc
					limit ` + pNValue

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TSC01) ", lErr)
		return &lTopNProdArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lTopNProdRec.Category, &lTopNProdRec.Qty)
			if lErr != nil {
				log.Println("Error (TSC02) ", lErr)
				return &lTopNProdArr, lErr
			} else {
				lTopNProdArr = append(lTopNProdArr, lTopNProdRec)
			}
		}
	}
	log.Println("TopSaleCategory - ")
	return &lTopNProdArr, nil
}

func TopSaleRegion(r *http.Request, pNValue string) (*[]models.TOPNPRODUCTS, error) {
	log.Println("TopSaleRegion + ")
	var lTopNProdRec models.TOPNPRODUCTS
	var lTopNProdArr []models.TOPNPRODUCTS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
					select *
					from (select o.region, sum(quantity) quantity
					from order_items oi , orders o 
					where oi.order_id = o.id and o.date_of_sale BETWEEN $1 AND $2
					group by o.region) t
					order by t.quantity desc
					limit ` + pNValue

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TSR01) ", lErr)
		return &lTopNProdArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lTopNProdRec.Region, &lTopNProdRec.Qty)
			if lErr != nil {
				log.Println("Error (TSR02) ", lErr)
				return &lTopNProdArr, lErr
			} else {
				lTopNProdArr = append(lTopNProdArr, lTopNProdRec)
			}
		}
	}
	log.Println("TopSaleRegion - ")
	return &lTopNProdArr, nil
}
