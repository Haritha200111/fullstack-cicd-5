package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func getDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// default for local docker-compose (service name "db")
		dsn = "postgresql://postgres:postgres@db:5432/mydb?sslmode=disable"
	}
	return sql.Open("postgres", dsn)
}

func main() {
	db, err := getDB()
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	// Wait for DB ready
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Printf("waiting for db... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/users", func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, email, "createdAt" FROM "User" ORDER BY id`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, u)
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.POST("/users", func(c *gin.Context) {
		var payload struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		var id int
		err := db.QueryRow(`INSERT INTO "User"(email) VALUES($1) RETURNING id`, payload.Email).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Backend running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
