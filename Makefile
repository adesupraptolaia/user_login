.PHONY: test
test:
	go test -v ./...

.PHONY: docker-build
docker-build:
	docker build -t user_login:latest .

.PHONY: build
build:
	go build -ldflags="-s -w" -o app

.PHONY: run
run:
	make docker-build
	docker-compose up -d

.PHONY: down
down:
	docker-compose down
# AUTH APP
.PHONY: run-auth
run-auth:
	go run main.go auth

.PHONY: docker-run-auth
docker-run-auth:
	docker run -e APP_NAME=auth --name auth_app user_login_app:latest

# User APP
.PHONY: run-user
run-user:
	go run main.go user

.PHONY: docker-run-user
docker-run-user:
	docker run -e APP_NAME=user --name user_app user_login_app:latest