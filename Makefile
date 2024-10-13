.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	go build cmd/main.go && ./main
