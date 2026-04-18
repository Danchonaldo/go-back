DB_URL=postgres://postgres:1234@localhost:5432/taskdb?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" force 1

migrate-version:
	migrate -path db/migrations -database "$(DB_URL)" version