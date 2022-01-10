package main

import (
	"pizzeria/api"
	"pizzeria/db"
	"pizzeria/repo"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := db.InitDB()
	if err != nil {
		logger.Fatal("InitDB", zap.Error(err))
	}

	s := repo.NewStorageSQL(db)
	o := repo.NewOrderSQL(db)
	srv := api.NewAPIService(logger.Named("api"), s, o)

	srv.Run()
}
