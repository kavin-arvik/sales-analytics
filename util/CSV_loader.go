package util

import (
	db "SalesAnalytics/DB"
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func LoadCSV(filePath string) error {
	var lErr error

	file, lErr := os.Open(filePath)
	if lErr != nil {
		return lErr
	}
	defer file.Close()
	var lResult sql.Result
	var lRowsAffected int64

	reader := csv.NewReader(file)
	records, lErr := reader.ReadAll()
	if lErr != nil {
		return lErr
	}

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		orderID := uuid.New()
		customerID := uuid.New()
		productID := uuid.New()

		name := row[12]
		email := row[13]
		address := row[14]

		lCustomerInsert := `
		    INSERT INTO customers (id, name, email, address)
		    VALUES ($1, $2, $3, $4)
		    ON CONFLICT (email) DO NOTHING
		`

		lResult, lErr = db.GDBConnection.Exec(lCustomerInsert, customerID, name, email, address)
		if lErr != nil {
			log.Println("Error : ULC01 ", lErr.Error())
			return lErr
		}
		lRowsAffected, _ = lResult.RowsAffected()

		//If the customer already exists, no insert occurs. So that, fetching the email's customers for orders table
		if lRowsAffected == 0 {
			customerID, lErr = selectIds("customers", "email", email)
			if lErr != nil {
				log.Println("Error : ULC02 ", lErr.Error())
				return lErr
			}

		}
		//--------------------------------------------------------------------------------------------------------------------------------------

		productName := row[3]
		category := row[4]

		lProductInsert := `
            INSERT INTO products (id, name, category)
            VALUES ($1, $2, $3)
            ON CONFLICT (name) DO NOTHING
        `
		lResult, lErr = db.GDBConnection.Exec(lProductInsert, productID, productName, category)
		if lErr != nil {
			log.Println("Error : ULC02 ", lErr.Error())
			return lErr
		}
		lRowsAffected, _ = lResult.RowsAffected()
		if lRowsAffected == 0 {
			productID, lErr = selectIds("products", "name", productName)
			if lErr != nil {
				log.Println("Error : ULC02 ", lErr.Error())
				return lErr
			}
		}

		//--------------------------------------------------------------------------------------------------------------------------------------

		saleDate, _ := time.Parse("2006-01-02", row[6])
		region := row[5]
		paymentMethod := row[11]

		lOrderInsert := `
            INSERT INTO orders (id, customer_id, date_of_sale, region, payment_method)
            VALUES ($1, $2, $3, $4, $5)
        `
		_, lErr = db.GDBConnection.Exec(lOrderInsert, orderID, customerID, saleDate, region, paymentMethod)
		if lErr != nil {
			log.Println("Error : ULC03 ", lErr.Error())
			return lErr
		}
		// lRowsAffected, _ = lResult.RowsAffected()

		//--------------------------------------------------------------------------------------------------------------------------------------

		qty, _ := strconv.Atoi(row[7])
		unitPrice, _ := strconv.ParseFloat(row[8], 64)
		discount, _ := strconv.ParseFloat(row[9], 64)
		shipping, _ := strconv.ParseFloat(row[10], 64)

		lOrderItemInsert := `
            INSERT INTO order_items (order_id, product_id, quantity, unit_price, discount, shipping_cost)
            VALUES ($1, $2, $3, $4, $5, $6)`

		_, lErr = db.GDBConnection.Exec(lOrderItemInsert, orderID, productID, qty, unitPrice, discount, shipping)
		if lErr != nil {
			log.Println("Error : ULC04 ", lErr.Error())
			return lErr
		}
		// lRowsAffected, _ = lResult.RowsAffected()
	}

	return nil
}

func selectIds(pTab, pCol, pDetails string) (lId uuid.UUID, lErr error) {
	log.Println("selectIds + ")
	lSelectString := `select id from ` + pTab + ` where ` + pCol + ` = $1`

	lRows, lErr := db.GDBConnection.Query(lSelectString, pDetails)
	if lErr != nil {
		log.Println("Error : ULC01.1 ", lErr.Error())
		return
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr = lRows.Scan(&lId)
			if lErr != nil {
				log.Println("Error : ULC01.2 ", lErr.Error())
				return
			}
		}
	}
	log.Println("selectIds - ")
	return
}
