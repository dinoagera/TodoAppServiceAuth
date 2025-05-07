include .env
MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"
migration-status:
	goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) status -v
migration-up:
	goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) up -v
migration-down:
	goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) down -v