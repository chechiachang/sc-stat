run:
	godotenv -f .env.local -- go run cmd/daemon/main.go

goose:
	goose -dir migrations sqlite3 ./db/development up
	goose -dir migrations sqlite3 ./db/development status
