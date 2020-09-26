package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/tkrouty/avitojob-trainee-task/models"
)

type DBWrapper struct {
	Conn *pgx.Conn
}

func (db *DBWrapper) MakeTransaction(t models.Transaction) error {
	if t.SourceID != "" {
		command := balanceEditStatement(t.SourceID, -t.Sum)
        fmt.Println(command)
		_, err := db.Conn.Exec(context.Background(), command)

		if err != nil {
			return err
		}

	}

	if t.TargetID != "" {
		command := balanceEditStatement(t.TargetID, t.Sum)
        fmt.Println(command)
		_, err := db.Conn.Exec(context.Background(), command)
		if err != nil {
			return err
		}
	}

    command := writeHistoryStatement(t.SourceID, t.TargetID, t.Sum, t.TransactionTime)
    fmt.Println(command)
    _, err := db.Conn.Exec(context.Background(), command)

    if err != nil {
        return err
    }

	return nil
}

func (db *DBWrapper) ShowBalance(u *models.User) error {
    command := showBalanceStatement(u.UserID)

    err := db.Conn.QueryRow(context.Background(), command).Scan(&u.Balance)

    if err != nil {
		return err
	}

    return nil
}

func (db *DBWrapper) ShowHistory(u models.User) ([]models.Transaction, error) {
    var res []models.Transaction

    command := showHistoryStatement(u.UserID)

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
