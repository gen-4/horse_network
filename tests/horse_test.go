package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/api/models"
)

func TestCreateHorse(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Name:  "Frederik",
		Age:   12,
		Breed: "PRI",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	if response.Data == nil {
		t.Errorf("Expected horse data, got nil")
	}
}

func TestCreateHorseNoName(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Age:   12,
		Breed: "PRI",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestCreateHorseNameTooShort(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Name:  "Fe",
		Age:   12,
		Breed: "PRI",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestCreateHorseIncorrectGender(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Name:   "Frederik",
		Age:    1,
		Gender: "fem",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestGetHorses(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Name:   "Frederik",
		Age:    1,
		Gender: "m",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	req, _ = http.NewRequest("GET", "/horses", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	if response.Data == nil {
		t.Errorf("Expected horse data, got nil")
	}

	if len(response.Data.([]any)) != 1 {
		t.Errorf("Expected value %d, got %d", 1, len(response.Data.([]models.Horse)))
	}
}

func TestDeleteHorse(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Name:   "Frederik",
		Age:    1,
		Gender: "m",
	}

	var horseId uint
	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)
	horseId = uint(response.Data.(map[string]any)["id"].(float64))

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/horse/%d", horseId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestDeleteHorseNotFoundHorse(t *testing.T) {
	SetupTestDB()
	user := AddUser(
		"Albert",
		25,
		"US",
		"albert@example.com",
		"m",
		"albert_pass",
		[]models.Role{GetRoleByName(models.USER)},
	)

	horse := models.Horse{
		Name:   "Frederik",
		Age:    1,
		Gender: "m",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(horse)
	req, _ := http.NewRequest("POST", "/horse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/horse/%d", UNEXISTENT_ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}
