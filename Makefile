.PHONY: dev
dev:
	go run cmd/gosynchro/main.go --verbose proxy --remote http://localhost:8080

.PHONY: build
build:
	go build -o bin/gosynchro cmd/gosynchro/main.go

.PHONY: run
run:
	./bin/gosynchro --verbose proxy --remote http://localhost:8080