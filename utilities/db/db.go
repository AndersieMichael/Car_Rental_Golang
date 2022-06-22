package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	
	_ "github.com/lib/pq"
)

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
)

func Connect() *sqlx.DB {
	host = "35.187.248.198"
	port = 5432
	user = "postgres"
	password = "d3v3l0p8015"
	dbname = "Car_Rental_Test_2"

	var dbInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("masuk 1")
		fmt.Println(err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.SetMaxOpenConns(10)                 //max
	db.SetConnMaxIdleTime(2 * time.Second) //max time connection may be idle
	db.SetConnMaxLifetime(5 * time.Second) // max time connection may be reuse

	return db
}
