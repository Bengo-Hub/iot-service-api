package main

import (
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/bengobox/iot-service/internal/config"
	"github.com/bengobox/iot-service/internal/ent"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Bypass PgBouncer for direct connection during seed.
	dbURL := cfg.Postgres.URL
	if cfg.Postgres.MigrateURL != "" {
		dbURL = cfg.Postgres.MigrateURL
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	_ = client // use client for seeding once Ent schema is finalized

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

	for _, r := range defaultRoles {
		log.Printf("Would seed role: %s - %s", r.code, r.name)
	}

	log.Println("seed completed")
}
