package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wolfalex23/go_final_project/pkg/api"
	"github.com/wolfalex23/go_final_project/pkg/db"
	"github.com/wolfalex23/go_final_project/pkg/server"
)

func main() {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)

	}

	defer db.Close()

	fmt.Println("Server running")
	api.Init()

	err = server.Run()
	if err != nil {
		panic(err)
	}

}
