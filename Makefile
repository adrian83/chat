
docker:
	sudo systemctl start docker

compose-up:
	sudo docker-compose up --build

deps:
	echo "starting rethinkdb image (version 2.4.0) and redis image (version 5.0.7)"
	docker run -p 28015:28015 -p 8080:8080 -d rethinkdb:2.4.0
	docker run -p 6379:6379 -d redis:5.0.7

be-lint:
	echo "running linters"
	golangci-lint run ./...

be-fmt:
	echo "formatting go code"
	gofmt -w .

be-test: 
	echo "running backend tests"
	go mod tidy
	go test ./... -cover


be-run: export SERVER_PORT=7070
be-run: export SERVER_HOST=0.0.0.0
be-run: export DATABASE_NAME=chat
be-run: export DATABASE_HOST=localhost
be-run: export DATABASE_PORT=28015
be-run: export STATICS_PATH=static
be-run: export SESSION_DB_NAME=0
be-run: export SESSION_DB_PASSWORD=
be-run: export SESSION_DB_HOST=localhost
be-run: export SESSION_DB_PORT=6379


be-run: 
	echo "running backend"
	go mod tidy
	go run cmd/chat/main.go

be-all: be-fmt be-test be-run

lint:
	echo "running linters" 
	golangci-lint run ./...