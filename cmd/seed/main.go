package main

import (
	"context"
	"log"

	"github.com/bengobox/iot-service/internal/config"
	"github.com/bengobox/iot-service/internal/platform/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	ctx := context.Background()
	client, err := database.NewClient(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer client.Close()

	// Ensure schema exists
	if err := database.RunMigrations(ctx, client); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Seed default roles
	defaultRoles := []struct {
		code        string
		name        string
		description string
		permissions []string
	}{
		{"admin", "Admin", "Full access to all IoT devices and settings", []string{"iot:devices:read", "iot:devices:write", "iot:devices:delete", "iot:devices:manage", "iot:telemetry:read"}},
		{"member", "Member", "Can manage assigned devices", []string{"iot:devices:read", "iot:devices:write", "iot:telemetry:read"}},
		{"viewer", "Viewer", "Read-only access to devices and telemetry", []string{"iot:devices:read", "iot:telemetry:read"}},
	}

	// Note: Role seeding will be implemented after Ent code generation
	// For now, just log that seeding is ready
	log.Println("Role seeding ready (Ent code must be generated first)")
	for _, r := range defaultRoles {
		log.Printf("Would seed role: %s - %s", r.code, r.name)
	}

	log.Println("seed completed")
}

