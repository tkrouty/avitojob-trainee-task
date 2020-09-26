package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tkrouty/avitojob-trainee-task/db"
	"github.com/tkrouty/avitojob-trainee-task/api"
)

func SetupRouter(db db.DBWrapper) *gin.Engine {
	r := gin.Default()

	r.HandleMethodNotAllowed = true

	accountAPI := api.AccountAPI{DB: db}
	// list of all routes
	r.POST("/add/:UserID", accountAPI.Add)
	r.POST("/deduct/:UserID", accountAPI.Deduct)
	r.POST("/transfer", accountAPI.Transfer)
	r.GET("/show_balance/:UserID", accountAPI.ShowBalance)
	r.GET("/show_history/:UserID", accountAPI.ShowHistory)

	return r
}
