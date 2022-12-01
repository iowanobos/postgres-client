LOCAL_BIN?=$(CURDIR)/bin

GOPATH?=$(HOME)/go
export GOPATH
GOBIN?=$(GOPATH)/bin
export GOBIN

# migrations
DB_DRIVER=postgres
MIGRATION_FORMAT=sql
MIGRATE_TAG=v4.14.1
MIGRATIONS_DIR=migrations
MONGO_URI?=postgres://root:root@localhost:5432/test?sslmode=disable

# migrate
MIGRATE_BIN=$(GOBIN)/migrate
$(MIGRATE_BIN):
	go install -tags '$(DB_DRIVER)' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_TAG)

# commands
.PHONY: migrate-create
migrate-create: $(MIGRATE_BIN)
	$(MIGRATE_BIN) create -dir $(MIGRATIONS_DIR) -ext $(MIGRATION_FORMAT) -seq $(name)

.PHONY: migrate-up
migrate-up: $(MIGRATE_BIN)
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(MONGO_URI)" -verbose up

.PHONY: migrate-down
migrate-down: $(MIGRATE_BIN)
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(MONGO_URI)" -verbose down 1

.PHONY: migrate-reset
migrate-reset: $(MIGRATE_BIN)
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(MONGO_URI)" -verbose force $(version)
