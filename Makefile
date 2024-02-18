
.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: serve
serve:
	@go run main.go

.PHONY: test
test:
	@cd ./test/handlers && go test -v

.PHONY: coverage
coverage:
	go test ./test/... -coverpkg ./...