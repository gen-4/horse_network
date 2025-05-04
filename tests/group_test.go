package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/api/handlers"
	"api/api/models"
)

func TestCreateGroup(t *testing.T) {
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

	group := models.Group{
		Name: "Montana WH Group",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)
}

func TestCreateGroupSameName(t *testing.T) {
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

	group := models.Group{
		Name: "Montana WH Group",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	req, _ = http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
}

func TestCreateGroupNoName(t *testing.T) {
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

	group := models.Group{}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestJoinGroup(t *testing.T) {
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

	group := models.Group{
		Name: "Montana WH Group",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	handlers.DB.First(&group)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/group/%d/join", group.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestJoinGroupNotFound(t *testing.T) {
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

	group := models.Group{
		Name: "Montana WH Group",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	handlers.DB.First(&group)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/group/%d/join", UNEXISTENT_ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}

func TestLeaveGroup(t *testing.T) {
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

	group := models.Group{
		Name: "Montana WH Group",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	handlers.DB.First(&group)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/group/%d/leave", group.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestLeaveGroupNotFound(t *testing.T) {
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

	group := models.Group{
		Name: "Montana WH Group",
	}

	router := GetRoutes()
	token := LogUser(router, user)
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/group", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	handlers.DB.First(&group)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/group/%d/leave", UNEXISTENT_ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}
