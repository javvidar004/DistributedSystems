package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type PasswordUpdateRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteUserRequest struct {
	Username string `json:"username"`
}

type UserModel struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func getUsers(db *sql.DB, c *gin.Context) {
	users, err := getUsersDB(db, c)
	var userResponses []UserResponse
	for user := range users {
		var userResponse UserResponse
		userResponse.ID = users[user].ID
		userResponse.Name = users[user].Name
		userResponse.Username = users[user].Username
		userResponses = append(userResponses, userResponse)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userResponses)
}

func getUsersDB(db *sql.DB, c *gin.Context) ([]UserModel, error) {
	var users []UserModel
	rows, err := db.QueryContext(c.Request.Context(), "SELECT id, name, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserModel

		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func registerUser(db *sql.DB, c *gin.Context) {
	var user RegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registeredUser := registerUserDB(db, c, user)
	if registeredUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": registeredUser.Error()})
		return
	}

	c.JSON(http.StatusCreated)
}

func registerUserDB(db *sql.DB, c *gin.Context, user RegisterRequest) error {
	_, err := db.ExecContext(c.Request.Context(), "INSERT INTO users (name, username, password) VALUES ($1, $2, $3)", user.Name, user.Username, user.Password)
	return err
}

func userLogin(db *sql.DB, c *gin.Context) {
	var login loginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := getUserByUsername(db, c, login.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if user.Password != login.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	c.JSON(http.StatusOK)
}

func getUserByUsername(db *sql.DB, c *gin.Context, username string) (*UserModel, error) {
	var user UserModel
	err := db.QueryRowContext(c.Request.Context(), "SELECT id, name, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Name, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func updatePassword(db *sql.DB, c *gin.Context) {
	var updateRequest PasswordUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := getUserByUsername(db, c, updateRequest.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if user.Password == updateRequest.NewPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password must be different from the old password"})
		return
	}
	if user.Password != updateRequest.OldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = updatePasswordDB(db, c, updateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"method": "PUT"})
}

func updatePasswordDB(db *sql.DB, c *gin.Context, updateRequest PasswordUpdateRequest) error {
	_, err := db.ExecContext(c.Request.Context(), "UPDATE users SET password = $1 WHERE username = $2", updateRequest.NewPassword, updateRequest.Username)
	return err
}

func deleteUser(db *sql.DB, c *gin.Context) {
	var deleteRequest DeleteUserRequest
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := deleteUserDB(db, c, deleteRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK)
}

func deleteUserDB(db *sql.DB, c *gin.Context, deleteRequest DeleteUserRequest) error {
	_, err := db.ExecContext(c.Request.Context(), "DELETE FROM users WHERE username = $1", deleteRequest.Username)
	return err
}

func main() {

	cfg := pq.Config{
		Host:     "localhost",
		Port:     5432,
		Database: "sampledb",
		User:     "root",
		Password: "root",
		SSLMode:  "disable",
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create connection pool.
	db := sql.OpenDB(c)
	defer db.Close()

	// Make sure it works.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/users", getUsers(db, gin.Context{}))
	router.POST("/login", userLogin(db, gin.Context{}))
	router.POST("/register", registerUser(db, gin.Context{}))
	router.PUT("/update", updatePassword(db, gin.Context{}))
	router.DELETE("/delete", deleteUser(db, gin.Context{}))

	router.Run(":8080")
}
