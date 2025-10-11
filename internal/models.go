package internal

import (
	"encoding/json"
	"time"
)
type Role string

const (
	Customer Role = "customer"
	Admin    Role = "admin"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

type OrderStatus string

const (
	Created    OrderStatus = "created"
	Dispatched OrderStatus = "dispatched"
	InTransit  OrderStatus = "in_transit"
	Delivered  OrderStatus = "delivered"
	Cancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID          int
	CustomerID  int
	Items       json.RawMessage
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt *time.Time
}
type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
