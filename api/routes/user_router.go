package routes

import (
	"banking_application/api/controllers"
	"banking_application/api/database"
	"banking_application/api/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, sqlDb *sql.DB) {
	db := database.NewDBContext(sqlDb)
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(*userService)

	incomingRoutes.GET("/users", userController.GetUsers())
	incomingRoutes.GET("/user/:user_id", userController.GetUser())
	incomingRoutes.POST("/user", userController.CreateUser())
}
