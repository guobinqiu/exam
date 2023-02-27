build:
	go build -o exam main.go

start: build
	./exam

test:
	go test -v ./rate

.PHONY: test
