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

func SalesCalculations(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("r.Method  - ", r.Method)

	log.Println("SalesCalculations + ")

	var lResp models.RESPONSE
	if r.Method == "GET" {

		lResp.Status = "S"
		lResp.Msg = "Data fetched successfully!!"
		lVars := mux.Vars(r)
		lId := lVars["id"]
		var lErr error

		if lId == "profitMargin" {
			lResp.Calculations, lErr = ProfitMargin(r)
		} else if lId == "clv" {
			lResp.Calculations, lErr = CLV(r)
		} else if lId == "customerSeg" {
			lResp.Calculations, lErr = customerSeg(r)
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
	log.Println("SalesCalculations - ")
}

func ProfitMargin(r *http.Request) (*[]models.CALCUALTIONS, error) {
	log.Println("ProfitMargin + ")
	var lCalcs models.CALCUALTIONS
	var lCalcsArr []models.CALCUALTIONS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				SELECT 
				p.name,
				ROUND(
					(
						SUM((oi.unit_price * oi.quantity) - (oi.discount + oi.shipping_cost)) 
						/ NULLIF(SUM(oi.unit_price * oi.quantity), 0)
					) * 100, 
					2
				) AS profit_margin_percentage
			FROM 
				orders o, order_items oi, products p
			WHERE o.id = oi.order_id and p.id = oi.product_id
				and o.date_of_sale BETWEEN $1 AND $2
			GROUP BY 
				oi.product_id, p."name";`

	log.Println("lSelectString - ", lSelectString)

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TC01) ", lErr)
		return &lCalcsArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lCalcs.Name, &lCalcs.ProfitMargin)
			if lErr != nil {
				log.Println("Error (TC02) ", lErr)
				return &lCalcsArr, lErr
			} else {
				lCalcsArr = append(lCalcsArr, lCalcs)
			}
		}
	}
	log.Println("ProfitMargin - ")
	return &lCalcsArr, nil
}

func CLV(r *http.Request) (*[]models.CALCUALTIONS, error) {
	log.Println("clv + ")
	var lCalcs models.CALCUALTIONS
	var lCalcsArr []models.CALCUALTIONS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				 select max("name") name, (avg(oi.quantity * oi.unit_price)) * case when (max(date_of_sale) - min(date_of_sale)) / 30 = 0 then 1
        													else (max(date_of_sale) - min(date_of_sale)) / 30
        													end as Average 
        from customers c , orders o, order_items oi 
        where c.id = o.customer_id and o.id = oi.order_id 
        and o.date_of_sale BETWEEN $1 AND $2
        group by email `

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TC01) ", lErr)
		log.Println("lSelectString - ", lSelectString)

		return &lCalcsArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lCalcs.Name, &lCalcs.CLV)
			if lErr != nil {
				log.Println("Error (TC02) ", lErr)
				return &lCalcsArr, lErr
			} else {
				lCalcsArr = append(lCalcsArr, lCalcs)
			}
		}
	}
	log.Println("clv - ")
	return &lCalcsArr, nil
}

func customerSeg(r *http.Request) (*[]models.CALCUALTIONS, error) {
	log.Println("customerSeg + ")
	var lCalcs models.CALCUALTIONS
	var lCalcsArr []models.CALCUALTIONS

	startDate, endDate := handlers.GetDateRange(r)

	lSelectString := `
				        SELECT 
						c.name,
						COUNT(DISTINCT o.id) AS total_orders,
						SUM(oi.quantity * oi.unit_price) AS total_spent,
						ROUND(AVG(oi.quantity * oi.unit_price), 2) AS avg_order_value,
						MAX(o.date_of_sale) AS last_order_date
					FROM customers c
					JOIN orders o ON c.id = o.customer_id
					JOIN order_items oi ON o.id = oi.order_id
					WHERE o.date_of_sale BETWEEN $1 AND $2
					GROUP BY c.id, c.name
					ORDER BY total_spent DESC;`

	lRows, lErr := db.GDBConnection.Query(lSelectString, startDate, endDate)
	if lErr != nil {
		log.Println("Error (TC01) ", lErr)
		return &lCalcsArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lCalcs.Name, &lCalcs.Total_orders, &lCalcs.Total_spent, &lCalcs.Avg_order_value, &lCalcs.Last_order_date)
			if lErr != nil {
				log.Println("Error (TC02) ", lErr)
				return &lCalcsArr, lErr
			} else {
				lCalcsArr = append(lCalcsArr, lCalcs)
			}
		}
	}
	log.Println("customerSeg - ")
	return &lCalcsArr, nil
}
