# CHAT-GO
Simple chat application written in [Go](https://golang.org/) and [Dart](https://www.dartlang.org/) programming languages with [Redis](//http://redis.io) as a session storage and [RethinkDB](https://www.rethinkdb.com/) as a main data storage.


### PREREQUISITES
1. Go
2. Dart
3. Docker


### RUNNING
Shortest path to run this application is:
1. `./run.sh dart`
2. `./run.sh run-infra`
3. `./run.sh go-deps`
4. `./run.sh run-go`


### QUALITY
1. No unit tests for now.
2. Checked with [metalinter](https://github.com/alecthomas/gometalinter)


### INFO FOR FUTURE ME
1. Write unit tests (Go and Dart)
2. Dart code refactoring
3. ...
