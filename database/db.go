package database

import (
	"os"

	"github.com/jinzhu/gorm"
	// Driver to connect to a Postgres database
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// Connection is the single connection instance
var Connection *gorm.DB

func init() {
	log.Info("Initializing database")
	var err error

	dbStringConnection := os.Getenv("DATABASE_URL")
	if dbStringConnection == "" {
		log.Fatal("No database string connection found. Set DATABASE_URL to continue.")
		return
	}

	Connection, err = gorm.Open("postgres", dbStringConnection)
	Connection.LogMode(true)

	if err != nil {
		log.Fatalf("Could not create a database connection - %s", err)
		return
	}
}

// Close closes de database connection
func Close() {
	log.Debug("About to close database connection...")
	err := Connection.Close()
	if err != nil {
		log.Warnf("Error while closing database. Please check for memory leak %s", err)
	} else {
		log.Info("DB Closed")
	}
}
