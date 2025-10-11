package internal

import (
	"context"
	"time"
)

// ===== Users =====
func (s *Store) CreateUser(ctx context.Context, u *User) error {
	now := time.Now()
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		u.Email, u.PasswordHash, u.Role, now,
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	query := `SELECT id,email,password_hash,role,created_at FROM users WHERE email=$1`
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ===== Orders =====
func (s *Store) CreateOrder(ctx context.Context, o *Order) error {
	now := time.Now()
	err := s.DB.QueryRowContext(ctx,
    `INSERT INTO orders (customer_id, items, status, created_at, updated_at)
     VALUES ($1, $2, $3, $4, $5) RETURNING id`,
    o.CustomerID, o.Items, o.Status, now, now,
).Scan(&o.ID)
if err != nil {
    return err
}
	o.CreatedAt = now
	o.UpdatedAt = now
	return nil
}

func (s *Store) GetOrderByID(ctx context.Context, id int) (*Order, error) {
	o := &Order{}
	query := `
		SELECT id, customer_id, items, status, created_at, updated_at, cancelled_at
		FROM orders
		WHERE id=$1
	`
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.CustomerID, &o.Items, &o.Status, &o.CreatedAt, &o.UpdatedAt, &o.CancelledAt,
	)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Store) CancelOrder(ctx context.Context, id int) error {
	now := time.Now()
	_, err := s.DB.ExecContext(ctx,
		"UPDATE orders SET status=$1, cancelled_at=$2, updated_at=$3 WHERE id=$4",
		Cancelled, now, now, id,
	)
	return err
}
