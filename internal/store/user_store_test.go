package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a temporary test database
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Create schema
	schema := `
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			created_at INTEGER NOT NULL
		);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestUserStore_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := NewUserStore(db)

	// Test creating user
	user, err := store.Create("alice")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.Username != "alice" {
		t.Errorf("expected username 'alice', got '%s'", user.Username)
	}

	if user.ID == "" {
		t.Error("expected non-empty user ID")
	}

	// Test duplicate username
	_, err = store.Create("alice")
	if err == nil {
		t.Error("expected error when creating duplicate user")
	}
}

func TestUserStore_GetByUsername(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := NewUserStore(db)

	// Create user
	created, err := store.Create("bob")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Get user
	retrieved, err := store.GetByUsername("bob")
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, retrieved.ID)
	}

	if retrieved.Username != "bob" {
		t.Errorf("expected username 'bob', got '%s'", retrieved.Username)
	}

	// Test non-existent user
	_, err = store.GetByUsername("nonexistent")
	if err == nil {
		t.Error("expected error when getting non-existent user")
	}
}
