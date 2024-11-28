go mod tidy
gofmt -w -s .
goimports -w .
go mod download
go vet ./...
golangci-lint run --color always
