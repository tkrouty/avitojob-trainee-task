package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tkrouty/avitojob-trainee-task/models"
)

type FinanceAPI struct {
	Manager FinanceManager
}

func (a *FinanceAPI) EditBalance(c *gin.Context) {
	t := models.Transaction{TransactionTime: time.Now()}
	if err := c.BindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if t.Sum > 0 {
		t.TargetID = c.Param("UserID")
	} else {
		t.SourceID = c.Param("UserID")
	}

	if err := a.Manager.makeTransaction(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "transaction completed"})
}

func (a *FinanceAPI) Transfer(c *gin.Context) {
	t := models.Transaction{TransactionTime: time.Now()}
	if err := c.BindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := a.Manager.makeTransaction(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "transaction is successful"})
}

func (a *FinanceAPI) ShowBalance(c *gin.Context) {
	u := models.User{UserID: c.Param("UserID")}

	currency := c.Query("currency")
	balance, err := a.Manager.getBalance(&u, currency)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if currency == "" {
		currency = "RUB"
	}

	c.JSON(200, gin.H{
		"user_id": u.UserID, "balance": balance, "currency": currency})
}

func (a *FinanceAPI) ShowHistory(c *gin.Context) {
	u := models.User{UserID: c.Param("UserID")}

	transactionHistory, err := a.Manager.getHistory(&u)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"user": u.UserID, "transaction_history": transactionHistory})
}
