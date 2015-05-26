//init.go
package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"os"
)

var (
	personRepo  *personRepository
	sessionRepo *sessionRepository
	sleepRepo   *sleepRepository
	dbmap       *gorp.DbMap
)

func GetDatabase() (*gorp.DbMap, error) {

	if dbmap == nil {

		//Get the database connection
		dbConnection := os.Getenv("DB_CONNECTION")

		//Try connecting to the database
		db, err := sql.Open("mysql", dbConnection)

		if err != nil {
			return nil, err
		}

		//Setup gorp if we have a dataase connection
		dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
		dbmap.AddTableWithName(Person{}, "person").SetKeys(true, "id")
		dbmap.AddTableWithName(Session{}, "session").SetKeys(true, "id")
		dbmap.AddTableWithName(Sleep{}, "sleep").SetKeys(true, "id")
	}

	return dbmap, nil
}
