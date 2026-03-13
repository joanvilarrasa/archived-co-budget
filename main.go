package main

import (
	"co-budget/data"
	"co-budget/lib"
	"co-budget/server"
	"context"
	"fmt"
)

func main() {
	db, dbErr := lib.OpenSQLite("./db.sqlite", "./sql")
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	data.InitAccountStore(db, context.Background(), dbErr)
	httpServer := server.NewHTTPServer()

	fmt.Println("listening on http://localhost:8080")
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
