package main

import (
	"context"
	"log"
	"os"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/tkrouty/avitojob-trainee-task/db"
	"github.com/tkrouty/avitojob-trainee-task/router"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// creating a logger
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	// connecting to DB
	logger.Println("connecting to", os.Getenv("POSTGRES_DB_URL"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_DB_URL"))
	if err != nil {
		logger.Printf("unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	logger.Println("connected to database")

	db := db.DBWrapper{Conn: conn, Logger: logger}

	// setting up router
	router := router.SetupRouter(db)
	endless.ListenAndServe(":8000", router)
}
