package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"

	"github.com/t3201v/ms/identity/internal/db/schema"
	"github.com/t3201v/ms/identity/internal/log"

	_ "github.com/lib/pq"
)

// IdentityDatabaseUrl is the environment variable for the identity database URL
const IdentityDatabaseUrl = "IDENTITY_DATABASE_URL"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if os.Getenv(IdentityDatabaseUrl) == "" {
		err := godotenv.Load()
		if err != nil {
			log.Error(ctx, "error loading .env file")
		}
	}
	databaseUrl := os.Getenv(IdentityDatabaseUrl)
	log.Config(log.LevelDebug, log.LevelDebug, os.Stdout)
	log.Debug(ctx, "database", "url", databaseUrl)

	if err := schema.Migrate(databaseUrl); err != nil {
		log.Error(ctx, "migrate", "error", err)
		return
	}
}
