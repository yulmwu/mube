.PHONY: tidy build run-apiserver run-mubelet

tidy:
	go mod tidy

build:
	go build ./...

run-apiserver:
	@if [ -f .env.apiserver ]; then \
		set -a; . ./.env.apiserver; set +a; \
	fi; \
	go run ./cmd/mube-apiserver

run-mubelet:
	@if [ -f .env.mubelet ]; then \
		set -a; . ./.env.mubelet; set +a; \
	fi; \
	go run ./cmd/mubelet
