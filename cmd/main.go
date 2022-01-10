package main

import (
	"log"
	"pizzeria/api"
	"pizzeria/db"
	"pizzeria/repo"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	s := repo.NewStorageSQL(db)
	o := repo.NewOrderSQL(db)
	srv := api.NewAPIService(s, o)

	srv.Run()
}
