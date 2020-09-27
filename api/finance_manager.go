package api

import (
	"errors"
	"sort"

	"github.com/tkrouty/avitojob-trainee-task/db"
	"github.com/tkrouty/avitojob-trainee-task/models"
)

type FinanceManager struct {
	DB db.DBWrapper
}

func (f *FinanceManager) makeTransaction(t *models.Transaction) error {
	if t.SourceID != "" {
		if balance, err := f.DB.GetBalance(t.SourceID); balance < t.Sum {
			if err != nil {
				return errors.New("transaction source is not registered")
			}
			return errors.New("transaction source has insufficient funds")
		}

		if err := f.DB.EditBalance(t.SourceID, -t.Sum); err != nil {
			return errors.New("transaction unsuccesful, database returned an error")
		}
	}

	if t.TargetID != "" {
		if err := f.DB.EditBalance(t.TargetID, t.Sum); err != nil {
			return errors.New("transaction unsuccesful, database returned an error")
		}
	}

	if err := f.DB.WriteHistory(t.SourceID, t.TargetID, t.Sum, t.TransactionTime); err != nil {
		return errors.New("couldn't write history, database returned an error")
	}

	return nil
}

func (f *FinanceManager) getBalance(u *models.User, currency string) (float64, error) {
	balance, err := f.DB.GetBalance(u.UserID)
	if err != nil {
		return 0, errors.New("unable to retrieve information about specified user_id")
	}

	if currency != "" {
		rate, err := getExchangeRate(currency)
		if err != nil {
			return 0, errors.New("unable to get exchange rate for " + currency)
		}
		balance *= rate
	}

	return balance, nil

}

func (f *FinanceManager) getHistory(u *models.User) ([]models.Transaction, error) {
	transaction_history, err := f.DB.GetHistory(u.UserID)
	if err != nil {
		return transaction_history, errors.New("unable to get transaction history, database returned an error")
	}
	sort.Slice(transaction_history, func(i, j int) bool {
		return transaction_history[j].TransactionTime.Before(transaction_history[i].TransactionTime)
	})
	sort.Slice(transaction_history, func(i, j int) bool {
		return transaction_history[j].Sum < transaction_history[i].Sum
	})

	return transaction_history, nil

}
