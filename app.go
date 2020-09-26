package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/tkrouty/avitojob-trainee-task/db"
	"github.com/tkrouty/avitojob-trainee-task/router"
)

func main() {
	// connecting to DB
	fmt.Println("connecting to", os.Getenv("POSTGRES_DB_URL"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_DB_URL"))
	fmt.Println("connected")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	db := db.DBWrapper{Conn: conn}

	// setting up router
	router := router.SetupRouter(db)
	router.Run()
}
