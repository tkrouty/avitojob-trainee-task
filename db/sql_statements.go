package db

import (
    "fmt"
    "time"
)

func balanceEditStatement(userID string, sum float64) string {
	return fmt.Sprintf(`
        UPDATE accounts SET balance=balance + %[2]v WHERE user_id='%[1]v';
        INSERT INTO accounts (user_id, balance) SELECT '%[1]v', %[2]v
		WHERE NOT EXISTS (SELECT 1 FROM accounts WHERE user_id='%[1]v');
        `, userID, sum)
}

func showBalanceStatement(userID string) string {
    return fmt.Sprintf(`SELECT balance FROM accounts WHERE user_id='%[1]v';`,
    userID)
}

func showHistoryStatement(userID string) string {
    return fmt.Sprintf(`SELECT * FROM transaction_history WHERE source_id='%[1]v'
        OR target_id='%[1]v';`, userID)
}

func writeHistoryStatement(sourceID string, targetID string,
    sum float64, transactionTime time.Time ) string {

    return fmt.Sprintf(`INSERT INTO transaction_history (source_id, target_id,
        sum, transaction_time) VALUES ('%[1]v', '%[2]v', %[3]v, '%[4]v');`,
    sourceID, targetID, sum, transactionTime.Format(time.RFC3339))
}
