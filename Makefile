# For local use
.PHONY: local-run
local-run:
	go run cmd/main.go

.PHONY: local-build
local-build:
	go build cmd/main.go && ./main

# For Docker
.PHONY: docker-db-up
docker-db-up:
	docker compose up -d go-storage-postgres

.PHONY: docker-db-stop
docker-db-stop:
	docker stop go-storage-postgres

.PHONY: docker-down
docker-down:
	docker compose down