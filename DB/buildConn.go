package db

import (
	"database/sql"
	"log"
)

var GDBConnection *sql.DB

func BuildConnecrtion() error {
	var lErr error
	GDBConnection, lErr = LocalDbConnect(Postgres)
	if lErr != nil {
		log.Println("Error in DB connect")
		return lErr
	}
	return nil
}
