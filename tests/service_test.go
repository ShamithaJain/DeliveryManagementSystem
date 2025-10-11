package tests

import (
    "context"
    "testing"
    "time"
    "encoding/json"
    "dms/internal"
)

func TestCreateAndCancelOrder(t *testing.T) {
    store, _ := internal.NewStore("postgres://postgres:MyDb!5432@localhost:5432/dms?sslmode=disable", "localhost:6379")
    
    items := json.RawMessage(`[{"name":"Laptop","quantity":1}]`)
    order := &internal.Order{
    CustomerID: 1,
    Items:      items,
    Status:     internal.Created,
    }


    if err := store.CreateOrder(context.Background(), order); err != nil {
        t.Fatalf("CreateOrder failed: %v", err)
    }

    // Start tracking in background
    go internal.StartOrderTracker(store, order.ID)

    // Cancel immediately
    if err := store.CancelOrder(context.Background(), order.ID); err != nil {
        t.Fatalf("CancelOrder failed: %v", err)
    }

    // Wait a second to ensure tracker has run
    time.Sleep(1 * time.Second)

    o, _ := store.GetOrderByID(context.Background(), order.ID)
    if o.Status != internal.Cancelled {
        t.Errorf("Expected Cancelled, got %s", o.Status)
    }
}
