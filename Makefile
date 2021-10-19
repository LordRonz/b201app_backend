generate-docs:
	@go run ./vendor/github.com/swaggo/swag/cmd/swag/main.go init -g pkg/api/api.go
