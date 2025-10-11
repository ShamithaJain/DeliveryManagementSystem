package internal

import (
	"context"
	"sync"
)

type OrderService struct {
	Store    *Store
	Trackers map[int]chan struct{}
	Mu       sync.RWMutex
}

// Constructor
func NewService(s *Store) *OrderService {
	return &OrderService{
		Store:    s,
		Trackers: make(map[int]chan struct{}),
	}
}

// Register a channel for an order
func (svc *OrderService) Register(orderID int) chan struct{} {
	svc.Mu.Lock()
	defer svc.Mu.Unlock()
	if ch, ok := svc.Trackers[orderID]; ok {
		return ch
	}
	ch := make(chan struct{})
	svc.Trackers[orderID] = ch
	return ch
}

// Unregister a channel after tracking
func (svc *OrderService) Unregister(orderID int) {
	svc.Mu.Lock()
	defer svc.Mu.Unlock()
	if ch, ok := svc.Trackers[orderID]; ok {
		close(ch)
		delete(svc.Trackers, orderID)
	}
}

// Cancel an order and notify tracker
func (svc *OrderService) CancelOrder(ctx context.Context, id int) error {
	_, err := svc.Store.DB.ExecContext(ctx,
		"UPDATE orders SET status=$1, cancelled_at=NOW() WHERE id=$2",
		"cancelled", id)
	if err != nil {
		return err
	}

	svc.Mu.RLock()
	if ch, ok := svc.Trackers[id]; ok {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
	svc.Mu.RUnlock()
	return nil
}
