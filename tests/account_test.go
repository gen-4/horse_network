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

func TestUpdateAccount(t *testing.T) {
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

	updateUser := models.UpdateUser{
		Username: "Ozuna",
		Country:  "PR",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(updateUser)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestUpdateAccountWrongUserName(t *testing.T) {
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

	updateUser := models.UpdateUser{
		Username: "Oz",
		Country:  "PR",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(updateUser)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestUpdateAccountWrongEmail(t *testing.T) {
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

	updateUser := models.UpdateUser{
		Mail: "ozuna.example.com",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(updateUser)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestUpdateAccountWrongGender(t *testing.T) {
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

	updateUser := models.UpdateUser{
		Gender: "male",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(updateUser)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestUpdateAccountWrongCountry(t *testing.T) {
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

	updateUser := models.UpdateUser{
		Country: "United Kingdom",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(updateUser)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}
