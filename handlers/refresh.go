package handlers

import (
	db "SalesAnalytics/DB"
	"SalesAnalytics/models"
	"SalesAnalytics/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func AutoRefresh() {
	for {
		now := time.Now()

		//refreshing and reloading the CSV file automaticlly everyday @ 6am
		if now.Hour() == 6 {
			lErr := RefreshSalesData("AUTOMATIC")
			if lErr != nil {
				log.Println("Error while reading - ", lErr.Error())
			}
		}
		//if not 6am, go to sleep for 30 mins
		time.Sleep(60 * 30 * time.Second)
	}
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var lResp models.RESPONSE
	lResp.Status = "S"
	lResp.Msg = "File loaded successfully!!"

	lErr := RefreshSalesData("MANUAL")
	if lErr != nil {
		log.Println("Error while reading - ", lErr.Error())
		lResp.Status = "E"
		lResp.Msg = "Error while reading"
	}

	lData, lErr := json.Marshal(lResp)

	if lErr != nil {
		fmt.Fprintf(w, "Error taking data"+lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}
}

func RefreshSalesData(refreshType string) error {
	log.Println("RefreshSalesData + ")
	start := time.Now()
	var logID int
	var lErr error
	var lSts, lErrorMsg string

	lInsertQry := `INSERT INTO refresh_logs (refresh_type, status, started_at) VALUES ($1, 'pending', $2) RETURNING id`

	lErr = db.GDBConnection.QueryRow(lInsertQry, refreshType, start).Scan(&logID)
	if lErr != nil {
		return lErr
	}

	// Actual CSV parsing
	lErr = util.LoadCSV("./data/sales_data.csv")

	end := time.Now()
	lSts = "SUCCESS"
	lUpdateString := `UPDATE refresh_logs SET status=$4, ended_at=$1, error_message=$2 WHERE id=$3`
	if lErr != nil {
		lSts = "FAILES"
		lErrorMsg = lErr.Error()
	}

	// Update log
	_, _ = db.GDBConnection.Exec(lUpdateString, end, lErrorMsg, logID, lSts)

	log.Println("RefreshSalesData - ")
	return lErr
}
