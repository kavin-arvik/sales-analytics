package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Structure to hold database connection details
type DatabaseType struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	DBType   string
	DB       string
}

// structure to hold all db connection details used in this program
type AllUsedDatabases struct {
	Postgres DatabaseType
}

// ---------------------------------------------------------------------------------
// function opens the db connection and return connection variable
// ---------------------------------------------------------------------------------
func LocalDbConnect(DBtype string) (*sql.DB, error) {
	DbDetails := new(AllUsedDatabases)
	DbDetails.Init()

	//log.Println(DbDetails)

	connString := ""
	localDBtype := ""

	var db *sql.DB
	var err error
	var dataBaseConnection DatabaseType
	log.Println("DBtype", DBtype)
	// get connection details
	if DBtype == DbDetails.Postgres.DB {
		dataBaseConnection = DbDetails.Postgres
		localDBtype = DbDetails.Postgres.DBType
	}

	// Prepare connection string
	if localDBtype == "postgres" {
		log.Println("IN", localDBtype)
		connString = `user=` + dataBaseConnection.User + ` password=` + dataBaseConnection.Password + ` dbname=` + dataBaseConnection.Database + ` host=` + dataBaseConnection.Server + ` sslmode=disable`
	}

	log.Println(localDBtype, "localDBtype")

	//make a connection to db
	if localDBtype != "" {
		db, err = sql.Open(localDBtype, connString)
		if err != nil {
			log.Println("Open connection failed:", err.Error())
		}
	} else {
		return db, fmt.Errorf(" Invalid DB Details")
	}

	return db, err
}

// --------------------------------------------------------------------
//
//	execute bulk inserts
//
// --------------------------------------------------------------------
func ExecuteBulkStatement(db *sql.DB, sqlStringValues string, sqlString string) error {
	log.Println("ExecuteBulkStatement+")
	//trim the last ,
	sqlStringValues = sqlStringValues[0 : len(sqlStringValues)-1]
	_, err := db.Exec(sqlString + sqlStringValues)
	if err != nil {
		log.Println(err)
		log.Println("ExecuteBulkStatement-")
		return err
	} else {
		log.Println("inserted Sucessfully")
	}
	log.Println("ExecuteBulkStatement-")
	return nil
}
