package db

import (
	"SalesAnalytics/common"
	"fmt"
	"strconv"
)

const (
	Postgres = "POSTGRES"
)

// Initializing DB Details
func (d *AllUsedDatabases) Init() {
	dbconfig := common.ReadTomlConfig("../dbconfig.toml")

	d.Postgres.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresServer"])
	d.Postgres.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresPort"]))
	d.Postgres.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresUser"])
	d.Postgres.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresPassword"])
	d.Postgres.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDatabase"])
	d.Postgres.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBType"])
	d.Postgres.DB = Postgres

}
