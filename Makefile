configure:
	gb vendor restore --all

build:
	gofmt -w src/waitinactivity
	go tool vet src/waitinactivity/*.go
	golint src/waitinactivity
	gb test
	gb build