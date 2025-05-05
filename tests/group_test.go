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
	group := AddGroup("Montana WH Group", user)

	router := GetRoutes()
	token := LogUser(router, user)

	handlers.DB.First(&group)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/group/%d/join", group.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
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
	group := AddGroup("Montana WH Group", user)

	router := GetRoutes()
	token := LogUser(router, user)

	handlers.DB.First(&group)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/group/%d/join", UNEXISTENT_ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
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
	group := AddGroup("Montana WH Group", user)

	router := GetRoutes()
	token := LogUser(router, user)

	handlers.DB.First(&group)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/group/%d/leave", group.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
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
	group := AddGroup("Montana WH Group", user)

	router := GetRoutes()
	token := LogUser(router, user)

	handlers.DB.First(&group)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/group/%d/leave", UNEXISTENT_ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}

func TestGetGroups(t *testing.T) {
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
	AddGroup("Montana WH Group", user)

	router := GetRoutes()
	token := LogUser(router, user)

	req, _ := http.NewRequest("GET", "/groups", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	groups := response.Data.([]any)
	if len(groups) != 1 {
		t.Errorf("Expected length %d, got %d", 1, len(groups))
	}
}

func TestGetGroupsJustOneBatch(t *testing.T) {
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
	for range handlers.GROUP_BATCH_SIZE + 1 {
		AddGroup("Montana WH Group", user)
	}

	router := GetRoutes()
	token := LogUser(router, user)

	req, _ := http.NewRequest("GET", "/groups", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	groups := response.Data.([]any)
	if len(groups) != handlers.GROUP_BATCH_SIZE {
		t.Errorf("Expected length %d, got %d", handlers.GROUP_BATCH_SIZE, len(groups))
	}
}

func TestGetGroupsIndex(t *testing.T) {
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
	const index = 2
	const indexName = "Z Index Group"
	for i := range handlers.GROUP_BATCH_SIZE + index {
		groupName := "Montana WH Group"
		if i == 0 { // Search is alphabetical so this should be placed last
			groupName = indexName
		}
		AddGroup(groupName, user)
	}

	router := GetRoutes()
	token := LogUser(router, user)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/groups?index=%d", index), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	groups := response.Data.([]any)
	if groups[handlers.GROUP_BATCH_SIZE-1].(map[string]any)["name"] != indexName {
		t.Errorf("Expected group name %s, got %s", indexName, groups[handlers.GROUP_BATCH_SIZE-1].(map[string]any)["name"])
	}
}

func TestGetGroupsBadRequest(t *testing.T) {
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
	for range handlers.GROUP_BATCH_SIZE + 1 {
		AddGroup("Montana WH Group", user)
	}

	router := GetRoutes()
	token := LogUser(router, user)

	req, _ := http.NewRequest("GET", "/groups?index=hello", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}
