#!/bin/bash

usage() {
	cat <<EOF

    Usage: $(basename $0) <command>

    run-docker            Starts systemd Docker daemon.
    run-rethink           Starts RethinkDB docker image.
    run-redis             Starts Redis docker image.
    run-infra             Starts Docker, RethinkDB and Redis
    run-go                Starts application (go run main.go)
    go-deps-get           Downloads Go dependencies
    go-deps-update        Update Go dependencies
    dart                  Downloads Dart dependencies and rebuilds Dart files
    dart-test             Runs Dart tests

EOF
	exit 1
}

run-docker() {
	set -e
		sudo systemctl start docker
	set +e
}

run-rethink() {
	set -e
		docker run -p 28016:28015 -p 8081:8080 -v $PWD/infrastructure/rethinkdb:/data -d rethinkdb
		echo "RethinkDB is listening on ports: 28016 and 8081. Data is stored inside 'infrastructure/rethinkdb' directory"
	set +e
}

run-redis() {
	set -e
		docker run -p 6380:6379 -v $PWD/infrastructure/redis:/data -d redis redis-server --appendonly yes
		echo "Redis is listening on port 6380. Data is stored inside 'infrastructure/redis' directory"
	set +e
}

run-infra() {
	set -e
		run-docker
		run-rethink
		run-redis
	set +e
}

dart() {
	set -e
		cd src/static && pub get && pub build
	set +e
}

dart-test() {
	set -e
		cd src/static && pub run test web/*test.dart
	set +e
}

run-go() {
	set -e
		cd src/ && go run main.go $PWD/config.json
	set +e
}

go-deps-update() {
	set -e
		cd src/ && go get -u ./...
	set +e
}

go-deps-get() {
	set -e
		cd src/ && go get ./...
	set +e
}

CMD="$1"
shift
case "$CMD" in
	run-docker)
		run-docker
	;;
	run-rethink)
		run-rethink
	;;
	run-redis)
		run-redis
	;;
	run-infra)
		run-infra
	;;
	run-go)
		run-go
	;;
	go-deps-update)
		go-deps-update
	;;
	go-deps-get)
		go-deps-get
	;;
	dart)
		dart
	;;
	dart-test)
		dart-test
	;;
	*)
		usage
	;;
esac
