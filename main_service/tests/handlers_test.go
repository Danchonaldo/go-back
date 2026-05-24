package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"main_service/db"
	"main_service/handlers"
	"main_service/middleware"
	"main_service/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	testDB.AutoMigrate(&models.User{}, &models.Board{}, &models.Task{})
	db.DB = testDB
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func getAuthToken(userID uint, email string) string {
	import_jwt := func() string { return "test_token" }
	_ = import_jwt
	return "Bearer test"
}

// Test 1: Register - success
func TestRegister_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/auth/register", handlers.Register)

	body := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User registered successfully", resp["message"])
}

// Test 2: Register - duplicate email
func TestRegister_DuplicateEmail(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/auth/register", handlers.Register)

	body := map[string]string{
		"name":     "Test User",
		"email":    "dup@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	// First registration
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)

	// Second registration with same email
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}

// Test 3: Register - missing fields
func TestRegister_MissingFields(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/auth/register", handlers.Register)

	body := map[string]string{"name": "No Email"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test 4: Login - wrong password
func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	regBody, _ := json.Marshal(map[string]string{
		"name": "User", "email": "login@test.com", "password": "correct123",
	})
	regReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), regReq)

	loginBody, _ := json.Marshal(map[string]string{
		"email": "login@test.com", "password": "wrongpass",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test 5: Login - success returns token
func TestLogin_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	regBody, _ := json.Marshal(map[string]string{
		"name": "Alice", "email": "alice@test.com", "password": "pass1234",
	})
	regReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), regReq)

	loginBody, _ := json.Marshal(map[string]string{
		"email": "alice@test.com", "password": "pass1234",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}

// Test 6: GetTasks - unauthorized
func TestGetTasks_Unauthorized(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/api/tasks", middleware.AuthMiddleware(), handlers.GetTasks)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tasks", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test 7: CreateBoard - missing title
func TestCreateBoard_MissingTitle(t *testing.T) {
	setupTestDB()
	r := setupRouter()

	r.POST("/api/boards", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.CreateBoard(c)
	})

	body, _ := json.Marshal(map[string]string{"description": "no title"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test 8: CreateBoard - success
func TestCreateBoard_Success(t *testing.T) {
	setupTestDB()

	// Need a user first
	user := models.User{Name: "Bob", Email: "bob@test.com", Password: "hashed"}
	db.DB.Create(&user)

	r := setupRouter()
	r.POST("/api/boards", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handlers.CreateBoard(c)
	})

	body, _ := json.Marshal(map[string]string{
		"title":       "My Board",
		"description": "Test board",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Board
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "My Board", resp.Title)
}

// Test 9: CreateTask - success
func TestCreateTask_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Name: "Carol", Email: "carol@test.com", Password: "hashed"}
	db.DB.Create(&user)
	board := models.Board{Title: "Board", UserID: user.ID}
	db.DB.Create(&board)

	r := setupRouter()
	r.POST("/api/tasks", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handlers.CreateTask(c)
	})

	body, _ := json.Marshal(map[string]interface{}{
		"title":    "Fix bug",
		"content":  "Fix the login bug",
		"board_id": board.ID,
		"priority": "high",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Task
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Fix bug", resp.Title)
	assert.Equal(t, "todo", resp.Status)
}

// Test 10: UpdateTaskStatus - invalid status
func TestUpdateTaskStatus_InvalidStatus(t *testing.T) {
	setupTestDB()

	user := models.User{Name: "Dave", Email: "dave@test.com", Password: "hashed"}
	db.DB.Create(&user)
	board := models.Board{Title: "Board", UserID: user.ID}
	db.DB.Create(&board)
	task := models.Task{Title: "Task", BoardID: board.ID, UserID: user.ID, Status: "todo"}
	db.DB.Create(&task)

	r := setupRouter()
	r.PATCH("/api/tasks/:id/status", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handlers.UpdateTaskStatus(c)
	})

	body, _ := json.Marshal(map[string]string{"status": "invalid-status"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/tasks/1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
