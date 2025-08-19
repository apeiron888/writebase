package test_utils

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupTestDatabase creates a test database instance
func SetupTestDatabase(t testing.TB) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create unique database name for test isolation
	dbName := "test_db_" + time.Now().Format("20060102150405")
	db := client.Database(dbName)
	t.Logf("Using test database: %s", dbName)

	return db, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := db.Drop(ctx); err != nil {
			t.Logf("Cleanup warning: %v", err)
		}
		if err := client.Disconnect(ctx); err != nil {
			t.Logf("Cleanup warning: %v", err)
		}
	}
}

// SetupBenchmarkDatabase creates a database for benchmarks
func SetupBenchmarkDatabase(b *testing.B) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		b.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return client.Database("bench_db")
}

// CleanupBenchmarkDatabase cleans up benchmark database
func CleanupBenchmarkDatabase(b *testing.B, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Drop(ctx); err != nil {
		b.Logf("Cleanup warning: %v", err)
	}
	if err := db.Client().Disconnect(ctx); err != nil {
		b.Logf("Cleanup warning: %v", err)
	}
}