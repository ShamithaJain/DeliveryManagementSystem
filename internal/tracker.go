package internal

import (
    "context"
    "time"
    "log"
)

// Simulate order lifecycle: created -> dispatched -> in_transit -> delivered
func StartOrderTracker(store *Store, orderID int) {
    statuses := []OrderStatus{Dispatched, InTransit, Delivered}
    for _, status := range statuses {
        time.Sleep(5 * time.Second) // wait 5 sec per stage (adjust as needed)
        ctx := context.Background()

        // Check if order was cancelled
        order, err := store.GetOrderByID(ctx, orderID)
        if err != nil {
            log.Println("Order not found:", err)
            return
        }
        if order.Status == Cancelled {
            log.Println("Order cancelled, stopping tracker for order:", orderID)
            return
        }

        // Update status
        _, err = store.DB.ExecContext(ctx,
            "UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2",
            status, orderID,
        )
        if err != nil {
            log.Println("Failed to update order status:", err)
            return
        }

        // Update Redis cache
        updatedOrder, _ := store.GetOrderByID(ctx, orderID)
        store.CacheOrder(ctx, updatedOrder)
        log.Println("Order", orderID, "updated to", status)
    }
}
