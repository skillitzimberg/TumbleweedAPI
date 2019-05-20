package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes a connection to the database.
func InitDB(dataSourceName string) {
	var err error
	fmt.Println("Initializing database connection . . .")

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dataSourceName)

		if err != nil {
			fmt.Printf("Unable to Open DB: %s... Retrying\n", err.Error())
			time.Sleep(time.Second * 2)
		} else if err = db.Ping(); err != nil {
			fmt.Printf("Unable to Ping DB: %s... Retrying\n", err.Error())
			time.Sleep(time.Second * 2)
		} else {
			err = nil
			break
		}
	}
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Connection to database successful!\n")
}
