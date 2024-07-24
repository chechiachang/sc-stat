
run:
	godotenv -f .env.local -- go run cmd/sc-stat/main.go

t:
	godotenv -f .env.local -- go run cmd/test/main.go

goose:
	goose -dir migrations sqlite3 ./db/development up
	goose -dir migrations sqlite3 ./db/development status

docker:
	docker build -t chechiachang/sc-stat:dev .
