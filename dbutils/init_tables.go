package dbutils

import (
	"database/sql"
	"log"
)

func Initialize(dbDriver *sql.DB) {
	statement, driverError := dbDriver.Prepare(image)
	if driverError != nil {
		log.Println(driverError)
	}
	// Create train table
	_, statementError := statement.Exec()
	if statementError != nil {
		log.Println("Table already exists!")
	}
	log.Println("All tables created/initialized successfully!")
}
