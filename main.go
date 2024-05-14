package main

import (
	"database/sql"
	"fmt"
	"os"

	"banking_application/api/database"
	"banking_application/api/routes"
	"banking_application/api/util"

	"github.com/gin-gonic/gin"
)

func main() {
	var inMemDb = false
	var db *sql.DB
	if inMemDb {

		// Initialize the database
		database.InitDatabase()

		uuid := util.GenerateUniqueAlphaNumeric(32)

		// Add entries to tables
		exist := database.CreateTable("Users", database.Entry{"ID": uuid, "FirstName": "Ayodeji", "LastName": "Bolanle", "PhoneNumber": "08180084896"})
		if !exist {
			fmt.Println("Error:", "Error occurred when creating table")
		}

		exist = database.CreateTable("Accounts", database.Entry{"ID": uuid, "AccountNumber": "8180084896", "UserId": 1})
		if !exist {
			fmt.Println("Error:", "Error occurred when creating table")
		}

		exist = database.CreateTable("Transactions", database.Entry{"ID": uuid, "Amount": 50.00, "AccountId": 1})
		if !exist {
			fmt.Println("Error:", "Error occurred when creating table")
		}
	} else {
		// Connect to the postgres database
		db = database.ConnectDB()
		defer db.Close()
		ctx := database.NewDBContext(db)

		migrationErr := database.Migration(ctx, "./api/database/db_schema.sql")
		if migrationErr != nil {
			fmt.Print("Error:", "Error occurred while migrating initial data", migrationErr)
			return
		}

	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router, db)
	routes.AccountRoutes(router, db)
	routes.TransactionRoutes(router, db)

	router.Run(":" + port)
}
