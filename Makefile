.PHONY: run build clean migrate-up migrate-down

BINARY_NAME=agent-care-tg

include .env
export 
run:
	go run main.go

build:
	go build -o $(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME) main

# Placeholder for future migration tool (e.g., golang-migrate)
migrate-up:
	goose -dir migrations postgres "${DATABASE_URL}" up

migrate-down:
	goose -dir migrations postgres "${DATABASE_URL}" down

migrate-status:
	goose -dir migrations postgres "${DATABASE_URL}" status

migrate-create:
	goose -dir migrations create ${name} sql

