package routes

import (
	"banking_application/api/controllers"
	"banking_application/api/database"
	"banking_application/api/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func AccountRoutes(incomingRoutes *gin.Engine, sqlDb *sql.DB) {
	db := database.NewDBContext(sqlDb)
	validator := services.NewTransactionValidator(db)
	accountService := services.NewAccountService(db, validator)
	accountController := controllers.NewAccountController(accountService)

	incomingRoutes.GET("/accounts", accountController.GetAccounts())
	incomingRoutes.GET("/account/:account_number", accountController.GetAccount())
	incomingRoutes.POST("/account/deactivate/:account_number", accountController.DeactivateAccount())
	incomingRoutes.POST("/account", accountController.CreateAccount())
}
