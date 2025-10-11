package tests

import (
	"encoding/json"
	"testing"
	"context"
	"dms/internal"
)

func TestCreateOrder(t *testing.T) {
    store, err := internal.NewTestStore()
    if err != nil {
        t.Fatalf("failed to create test store: %v", err)
    }
    defer store.Close()

    order := &internal.Order{
    CustomerID: 1,
    Items:      json.RawMessage(`[{"name":"Laptop","quantity":1}]`),
    Status:     "created",
}

// Inside CreateOrder, do:

    if err := store.CreateOrder(context.Background(), order); err != nil {
        t.Fatalf("Failed to create order: %v", err)
    }

    if order.ID == 0 {
        t.Fatalf("Order ID not set")
    }
}

