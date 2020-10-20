# CHAT-GO
Simple chat application written in [Go](https://golang.org/) and [Dart](https://www.dartlang.org/) programming languages with [Redis](//http://redis.io) as a session storage and [RethinkDB](https://www.rethinkdb.com/) as a main data storage.

## Running

### Running with docker compose

#### Prerequisites
- Docker
- Docker Compose

#### Steps
1. Run `docker-compose up`
2. Navigate in browser to `localhost:7070`

### Running locally

#### Prerequisites
- Docker
- Go
- Dart (with webdev)
- Gnu Make

#### Steps
1. Start Infrastructure (RethinkDB and Redis): `make deps`
2. Build frontend `make fe-all`
3. Start backend: `make be-all`
4. Navigate in browser to `localhost:7070`
