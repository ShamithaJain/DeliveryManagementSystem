package tests

import (
	"context"
	"testing"

	"dms/internal"
)

func TestCreateAndGetUser(t *testing.T) {
	store, err := internal.NewTestStore()
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}
	defer store.Close()

	user := &internal.User{
		Email:        "test@example.com",
		PasswordHash: "password123",
		Role:         "customer",
	}

	if err := store.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	got, err := store.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	if got.Email != user.Email {
		t.Fatalf("Expected email %s, got %s", user.Email, got.Email)
	}
}
