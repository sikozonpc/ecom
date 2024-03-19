package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/sikozonpc/ecom/configs"
)

func NewMySQLStorage() (*sql.DB, error) {
	db, err := sql.Open("mysql", configs.Envs.GetDBConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
