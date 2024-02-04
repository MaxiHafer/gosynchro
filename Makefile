.PHONY: dev
dev:
	go run main.go --verbose proxy --remote http://localhost:8080

.PHONY: build
build:
	go build -o bin/gosynchro .

.PHONY: run
run: build
	./bin/gosynchro --verbose proxy --remote http://localhost:8080