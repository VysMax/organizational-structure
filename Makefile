#DB_DSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
DB_DSN = host=localhost port=5432 user=postgres password=password dbname=org_str sslmode=disable
MIGRATIONS_DIR = ./migrations

test-mocks:
	go test ./usecase

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

gen-mocks:
	mockgen -source=usecase/usecase.go \
	-destination=usecase/mocks/mock.go \
	-package=mocks