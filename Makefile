
run:
	godotenv -f .env.local -- go run cmd/sc-stat/main.go

t:
	godotenv -f .env.local -- go run cmd/test/main.go

goose:
	goose -dir migrations sqlite3 ./db/development up
	goose -dir migrations sqlite3 ./db/development status

image:
	docker buildx inspect
	docker buildx bake --print
	docker buildx bake --push

deploy:
	helm upgrade --install --reset-values --cleanup-on-fail --namespace sc-stat sc-stat charts/sc-stat
