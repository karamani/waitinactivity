configure:
	gb vendor update --all

build:
	gofmt -w src/waitinactivity
	go tool vet src/waitinactivity/*.go
	gb test
	gb build