package internal

import (
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// SetupRouter returns a router ready for integration tests
func SetupRouter(store *Store, jwtSecret []byte) *mux.Router {
	return NewRouter(store, jwtSecret)
}

// Helper to create a test server
func NewTestServer(store *Store, jwtSecret []byte) *httptest.Server {
	router := SetupRouter(store, jwtSecret)
	return httptest.NewServer(router)
}
