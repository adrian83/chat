
docker:
	sudo systemctl start docker

compose-build:
	sudo docker-compose build

compose-up:
	sudo docker-compose up

deps:
	echo "starting rethinkdb image (version 2.4.0) and redis image (version 5.0.7)"
	docker run -p 28015:28015 -p 8080:8080 -d rethinkdb:2.4.0
	docker run -p 6379:6379 -d redis:5.0.7

fe-get:
	echo "getting frontend dependencies" 
	cd static && pub get 

fe-build:
	echo "building frontend" 
	cd static && webdev build

fe-test:
	echo "running frontend tests" 
	cd static && pub run test web/*test.dart

fe-all: fe-get fe-build

be-lint:
	echo "running linters" 
	golangci-lint run ./...

be-fmt:
	echo "formatting go code"
	gofmt -w .

be-test: 
	echo "running backend tests"
	go test ./... -cover

be-run: 
	echo "running backend"
	go run cmd/chat/main.go config.json

be-all: be-fmt be-test be-run

lint:
	echo "running linters" 
	golangci-lint run ./...