package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// User represents a user model
type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phonenumber"`
}

// Initialize database
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		phonenumber TEXT NOT NULL UNIQUE
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

// Get all users
func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, username, phonenumber FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.PhoneNumber); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning user"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// Create a new user
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO users(username, phonenumber) VALUES(?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// Logger Middleware
	r.Use(gin.Logger())

	// CORS Middleware
	r.Use(cors.Default())

	// User routes
	api := r.Group("/api")
	{
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
	}

	log.Println("Server running on :8080")
	r.Run(":8080")
}
