package main

import (
	"database/sql"
	"log"
	"yildizskylab/src/api"
	"yildizskylab/src/db/sqlc"
	"yildizskylab/src/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(sqlc.New(db), config.Secret)
	server.Start(config.ServerAddress)
}

// TODO: CurrentUser

// TODO: Mailing system
// TODO: Swegger eklenecek

// SKY LOG
