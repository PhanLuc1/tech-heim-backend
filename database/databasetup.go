package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func DBset() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/ecommerce?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Susscessfully connected to MYSQL")
	return db
}

var Client *sql.DB = DBset()