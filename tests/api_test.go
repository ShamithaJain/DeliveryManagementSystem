package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"dms/internal"
)

var jwtSecret = []byte("testsecret")

func setupTestServer(t *testing.T) *httptest.Server {
	store, err := internal.NewTestStore()
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	router := internal.NewRouter(store, jwtSecret)
	ts := httptest.NewServer(router)
	t.Cleanup(func() {
		ts.Close()
		store.Close()
	})

	return ts
}

func TestSignupLogin(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Signup
	signupPayload := map[string]string{
		"email":         "integ@example.com",
		"password_hash": "pass123",
		"role":          "customer",
	}
	b, _ := json.Marshal(signupPayload)
	resp, err := http.Post(ts.URL+"/signup", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("Signup request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", resp.StatusCode)
	}

	// Login
	loginPayload := map[string]string{
		"email":         "integ@example.com",
		"password_hash": "pass123",
	}
	b2, _ := json.Marshal(loginPayload)
	resp2, err := http.Post(ts.URL+"/login", "application/json", bytes.NewReader(b2))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp2.StatusCode)
	}

	var data map[string]string
	if err := json.NewDecoder(resp2.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	if data["token"] == "" {
		t.Errorf("Expected token in response, got empty")
	}
}

func TestCreateOrderAndCancel(t *testing.T) {
	ts := setupTestServer(t)

	// Signup
	signupPayload := map[string]string{
		"email":         "orderuser@example.com",
		"password_hash": "pass123",
		"role":          "customer",
	}
	b, _ := json.Marshal(signupPayload)
	resp, err := http.Post(ts.URL+"/signup", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("Signup request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created for signup, got %d", resp.StatusCode)
	}

	// Login
	loginPayload := map[string]string{
		"email":         "orderuser@example.com",
		"password_hash": "pass123",
	}
	b2, _ := json.Marshal(loginPayload)
	respLogin, err := http.Post(ts.URL+"/login", "application/json", bytes.NewReader(b2))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer respLogin.Body.Close()
	if respLogin.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK for login, got %d", respLogin.StatusCode)
	}

	var loginResp map[string]string
	err = json.NewDecoder(respLogin.Body).Decode(&loginResp)
	if err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}
	token := loginResp["token"]

	// Create order
	orderPayload := map[string]interface{}{
		"items": []map[string]interface{}{
			{"name": "Laptop", "quantity": 1},
		},
	}
	b3, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", ts.URL+"/orders", bytes.NewReader(b3))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Create order request failed: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK for create order, got %d", resp2.StatusCode)
	}

	// Cancel order
	reqCancel, _ := http.NewRequest("POST", ts.URL+"/orders/1/cancel", nil)
	reqCancel.Header.Set("Authorization", "Bearer "+token)
	respCancel, err := http.DefaultClient.Do(reqCancel)
	if err != nil {
		t.Fatalf("Cancel order request failed: %v", err)
	}
	defer respCancel.Body.Close()
	if respCancel.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK for cancel order, got %d", respCancel.StatusCode)
	}
}
