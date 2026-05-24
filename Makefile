.PHONY: up down build test migrate logs

# Start all services
up:
	docker-compose up --build -d

# Stop all services
down:
	docker-compose down

# Build images
build:
	docker-compose build

# View logs
logs:
	docker-compose logs -f

# Run unit tests
test:
	cd main_service && go test ./tests/... -v -count=1

# Run migrations manually
migrate-up:
	cd main_service && \
	migrate -path db/migrations \
	        -database "postgres://postgres:postgres@localhost:5432/taskboard?sslmode=disable" \
	        up

migrate-down:
	cd main_service && \
	migrate -path db/migrations \
	        -database "postgres://postgres:postgres@localhost:5432/taskboard?sslmode=disable" \
	        down

# Run main service locally (needs DB running)
run-main:
	cd main_service && go run main.go

# Run notification service locally
run-notify:
	cd notification_service && go run main.go

# Run frontend locally
run-frontend:
	cd frontend && npm start

# Install frontend deps
install-frontend:
	cd frontend && npm install

# Clean docker
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f
