package api

import (
	"github.com/gin-gonic/gin"
    "github.com/tkrouty/avitojob-trainee-task/db"
    "github.com/tkrouty/avitojob-trainee-task/models"
    "time"
)

type AccountAPI struct {
	DB db.DBWrapper
}

func (a *AccountAPI) Add(c *gin.Context) {
	t := models.Transaction{TransactionTime: time.Now()}
	if err := c.BindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

    t.TargetID = c.Param("UserID")

    if err := a.DB.MakeTransaction(t); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
		return
    }

	c.JSON(200, gin.H{
		"message": "transaction completed"})
}

func (a *AccountAPI) Deduct(c *gin.Context) {
	t := models.Transaction{TransactionTime: time.Now()}
	if err := c.BindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

    t.SourceID = c.Param("UserID")

	if err := a.DB.MakeTransaction(t); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
		return
    }

    c.JSON(200, gin.H{
		"message": "transaction is successful"})
}

func (a *AccountAPI) Transfer(c *gin.Context) {
    t := models.Transaction{TransactionTime: time.Now()}
    if err := c.BindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

    if err := a.DB.MakeTransaction(t); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
		return
    }

    c.JSON(200, gin.H{
		"message": "transaction is successful"})
}

func (a *AccountAPI) ShowBalance(c *gin.Context) {
    u := models.User{}

    u.UserID = c.Param("UserID")

    if err := a.DB.ShowBalance(&u); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
		return
    }

	var currency string

	if currency = c.Query("currency"); currency != "" {
		rate, err := getExchangeRate(currency)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		u.Balance *= rate
	} else {
		currency = "RUB"
	}

    c.JSON(200, gin.H{
		"data": u, "currency" : currency})
}

func (a *AccountAPI) ShowHistory(c *gin.Context) {
    u := models.User{}

    u.UserID = c.Param("UserID")


    transactionHistory, err := a.DB.ShowHistory(u)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
		return
    }

    c.JSON(200, gin.H{"user": u, "transaction_history": transactionHistory})
}
