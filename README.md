API postman: https://dark-eclipse-55522.postman.co/workspace/Pyre~fa15a870-b5fc-4f96-a7f8-d5f65d20ccbc/collection/20536686-46c3db49-d794-4e99-b60e-6259265e181c?action=share

Migrate cmd: migrate -database "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable" -path core-internal/sqlc/migrations up

init file migrate: migrate create -ext sql -dir ./core/sqlc/migrations -seq init_tables
