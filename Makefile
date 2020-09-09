all: test

test:
	@goimports -w .
	@go test -timeout 10s -race -count 10 -cover -coverprofile=./envx.cover ./...

cover: test
	@go tool cover -html=./envx.cover
