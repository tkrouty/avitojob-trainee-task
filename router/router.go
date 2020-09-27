package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tkrouty/avitojob-trainee-task/api"
	"github.com/tkrouty/avitojob-trainee-task/db"
)

func SetupRouter(db db.DBWrapper) *gin.Engine {
	r := gin.Default()

	r.HandleMethodNotAllowed = true
	financeManager := api.FinanceManager{DB: db}
	financeAPI := api.FinanceAPI{Manager: financeManager}
	// list of all routes
	r.POST("/edit_balance/:UserID", financeAPI.EditBalance)
	r.POST("/transfer", financeAPI.Transfer)
	r.GET("/show_balance/:UserID", financeAPI.ShowBalance)
	r.GET("/show_history/:UserID", financeAPI.ShowHistory)

	return r
}
