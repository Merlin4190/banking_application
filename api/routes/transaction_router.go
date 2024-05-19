package routes

import (
	"banking_application/api/controllers"
	"banking_application/api/database"
	"banking_application/api/http"
	"banking_application/api/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(incomingRoutes *gin.Engine, sqlDb *sql.DB) {
	db := database.NewDBContext(sqlDb)
	client := http.NewClient()
	validator := services.NewTransactionValidator(db)
	transactionService := services.NewTransactionService(db, client, validator)
	transactionController := controllers.NewTransactionController(transactionService)

	//incomingRoutes.GET("/transactions", transactionController.GetTransactions())
	//incomingRoutes.GET("/transaction/:trans_id", transactionController.GetTransaction())
	incomingRoutes.POST("/transaction/deposit", transactionController.Deposit())
	incomingRoutes.POST("/transaction/withdraw", transactionController.Withdraw())
	incomingRoutes.POST("/transaction/transfer", transactionController.Transfer())
}
