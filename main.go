package main

import (
	"BreezeServer/api"
	db "BreezeServer/db/sqlc"
	"BreezeServer/util"
	"database/sql"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot laod config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
