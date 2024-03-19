package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sikozonpc/ecom/cmd/api"
	"github.com/sikozonpc/ecom/configs"
	"github.com/sikozonpc/ecom/db"
)

func main() {
	db, err := db.NewMySQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
