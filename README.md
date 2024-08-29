API postman: https://www.postman.com/dark-eclipse-55522/workspace/pyre/collection/20536686-46c3db49-d794-4e99-b60e-6259265e181c?action=share&creator=20536686

Migrate cmd: migrate -database "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable" -path core-internal/sqlc/migrations up

init file migrate: migrate create -ext sql -dir ./core/sqlc/migrations -seq init_tables