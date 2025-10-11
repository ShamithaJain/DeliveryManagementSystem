package internal

import (
    "github.com/gorilla/mux"
)
func NewRouter(store *Store, jwtSecret []byte) *mux.Router {
    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/signup", SignupHandler(store)).Methods("POST")
    r.HandleFunc("/login", LoginHandler(store, jwtSecret)).Methods("POST")

    // Authenticated routes
    auth := r.NewRoute().Subrouter()
    auth.Use(AuthMiddleware(jwtSecret))
    auth.HandleFunc("/orders", CreateOrderHandler(store)).Methods("POST")
    auth.HandleFunc("/orders/{id}/cancel", CancelOrderHandler(store)).Methods("POST")
    auth.HandleFunc("/admin/orders", GetAllOrdersHandler(store)).Methods("GET")
    auth.HandleFunc("/orders/{id}/track", TrackOrderHandler(store)).Methods("GET")

    return r
}
