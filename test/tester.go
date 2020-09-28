package test

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

type Tester struct {
	Logger *log.Logger
	DB     *pgx.Conn
}

func InitTester() Tester {
	tester := Tester{}
	file, err := os.OpenFile("test_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	tester.Logger = log.New(file, "TESTLOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return tester
}

func (t *Tester) ConnectToDB() {
	t.Logger.Println("connecting to", os.Getenv("TEST_DB_URL"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Logger.Printf("unable to connect to database: %v\n", err)
		return
	}
	t.Logger.Println("Connected")
	t.DB = conn
}

func (t *Tester) CloseDBConnection() {
	t.Logger.Println("Closing DB connection")
	t.DB.Close(context.Background())
}

func (t *Tester) FlushTables() {
	t.Logger.Println("flushing tables before a new test")
	_, err := t.DB.Exec(context.Background(), `DELETE FROM accounts; DELETE FROM transaction_history;`)
	if err != nil {
		t.Logger.Printf("errored while flushing tables %v\n", err)
	}
}
