package main

import (
	"currentPrice/database"
	"currentPrice/server"
	"time"
)

func main() {
	database.DBCon()
	for {
		server.Run()

		time.Sleep(1 * time.Second)
	}
}
