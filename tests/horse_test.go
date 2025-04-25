package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/api/handlers"
	"api/api/models"

	"github.com/gin-gonic/gin"
)

func TestCreateHorse(t *testing.T) {
	SetupTestDB()
	router := gin.Default()
	router.POST("/horse", handlers.CreateHorse)

	horse := models.Horse{}

	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(responseRecorder.Body).Decode(&response)

	if response.Data == nil {
		t.Errorf("Expected horse data, got nil")
	}
}
