package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/tkrouty/avitojob-trainee-task/models"
)

type DBWrapper struct {
	Conn   *pgx.Conn
	Logger *log.Logger
}

func (db *DBWrapper) EditBalance(userID string, sum float64) error {
	command := balanceEditStatement(userID, sum)
	db.Logger.Println(command)
	if _, err := db.Conn.Exec(context.Background(), command); err != nil {
		return err
	}

	return nil
}

func (db *DBWrapper) WriteHistory(sourceID string, targetID string,
	sum float64, transactionTime time.Time) error {

	command := writeHistoryStatement(sourceID, targetID, sum, transactionTime)
	db.Logger.Println(command)
	if _, err := db.Conn.Exec(context.Background(), command); err != nil {
		return err
	}

	return nil
}

func (db *DBWrapper) GetBalance(userID string) (float64, error) {
	var balance float64

	command := showBalanceStatement(userID)
	db.Logger.Println(command)

	if err := db.Conn.QueryRow(context.Background(), command).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}

func (db *DBWrapper) GetHistory(userID string) ([]models.Transaction, error) {
	var res []models.Transaction

	command := showHistoryStatement(userID)
	db.Logger.Println(command)
	rows, err := db.Conn.Query(context.Background(), command)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		err = rows.Scan(
			&t.TransactionID,
			&t.SourceID,
			&t.TargetID,
			&t.Sum,
			&t.TransactionTime,
		)
		if err != nil {
			return res, err
		}
		res = append(res, t)
	}

	return res, nil
}
