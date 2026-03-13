.PHONY: run build clean migrate-up migrate-down

BINARY_NAME=agent-care-tg

run:
	go run main.go

build:
	go build -o $(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME) main

# Placeholder for future migration tool (e.g., golang-migrate)
migrate-up:
	@echo "Running migrations up..."

migrate-down:
	@echo "Running migrations down..."
