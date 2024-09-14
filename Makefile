-include ./.env

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -cover -bench=. -benchmem -benchtime=100x ./...