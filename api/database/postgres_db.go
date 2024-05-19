package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"banking_application/api/util"
)

func ConnectDB() *sql.DB {

	/*host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")*/

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "@Merlino07"
		dbname   = "banking_application_db"
	)

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	// defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}
	fmt.Println("Connected to the database!")
	return db
}

func Migration(ctx *AppDBContext, filename string) error {
	// Read SQL commands from the file
	sqlCommands, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %v", err)
	}

	// Split SQL commands by semicolon
	commands := strings.Split(string(sqlCommands), ";")

	// Execute SQL commands
	for _, command := range commands {
		command = strings.TrimSpace(command)
		if command != "" {
			_, err := ctx.DB.Exec(command)
			if err != nil {
				log.Println("error here, ", err)
				return fmt.Errorf("error executing SQL command: %s\n%s\n", command, err)
			}
		}
	}

	// Seed initial data into tables
	if err := SeedData(ctx); err != nil {
		return fmt.Errorf("error seeding initial data: %v", err)
	}

	return nil
}

func SeedData(ctx *AppDBContext) error {

	// Check if data already exists in the users table
	var userCount int
	err := ctx.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return fmt.Errorf("error querying users table: %v", err)
	}

	// If userCount > 0, data already exists, so return without seeding
	if userCount > 0 {
		return nil
	}

	// Seed initial data into users table
	defaultPassword := "password123"

	// Hash the default password
	hashedPassword, err := util.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("password hashing failed")
	}

	if _, err := ctx.Insert(`INSERT INTO users (id, firstname, lastname, email, password, created_at) 
		VALUES ('4a604e74-ef0f-4f46-9b15-6bb24e3f2a06', 'Ayodeji', 'Bolanle', 'ayodeji_bolanle@yahoo.com', $1, $2)`, hashedPassword, time.Now()); err != nil {
		return fmt.Errorf("error seeding data into users table: %v", err)
	}

	// Seed initial data into accounts table
	if _, err := ctx.Insert(`INSERT INTO accounts (id, user_id, account_number, account_balance, is_active, created_at) 
		VALUES ('6c70ab16-4959-4286-b7a6-f1b219be091b', '4a604e74-ef0f-4f46-9b15-6bb24e3f2a06', '1234567890', 100.00, true, $1)`, time.Now()); err != nil {
		return fmt.Errorf("error seeding data into accounts table: %v", err)
	}

	// Seed initial data into transactions table
	if _, err := ctx.Insert(`INSERT INTO transactions (account_id, transaction_reference, transaction_type, transaction_record_type, 
                          transaction_amount, transaction_status, created_at) VALUES ('6c70ab16-4959-4286-b7a6-f1b219be091b', '123456789123456', 
                                                                                      'deposit', 'credit', 50.00, 'success', $1)`, time.Now()); err != nil {
		return fmt.Errorf("error seeding data into transactions table: %v", err)
	}

	return nil
}
