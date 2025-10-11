package tests

import (
	"context"
	"encoding/json"
	"testing"

	"dms/internal"
)

func TestCreateAndGetOrder(t *testing.T) {
	store, err := internal.NewTestStore()
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}
	defer store.Close()

	items := json.RawMessage(`[{"name":"Laptop","quantity":1}]`)
	order := &internal.Order{
		CustomerID: 1,
		Items:      items,
		Status:     internal.Created,
	}

	// Create order
	if err := store.CreateOrder(context.Background(), order); err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	// Get order
	o, err := store.GetOrderByID(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("GetOrderByID failed: %v", err)
	}

	if o.ID != order.ID {
		t.Fatalf("Expected order ID %d, got %d", order.ID, o.ID)
	}
}
