default: go-run
.PHONY: go-run

# Golang
go-run:
	go run cmd/main.go -config ./configLocal.env

go-run-dev:
	go run cmd/main.go -config ./configDev.env

go-fmt:
	go fmt ./...

go-lint:
	golangci-lint run ./... -c .golangci.yml

go-memory-check:
	fieldalignment ./...

go-memory-fix:
	fieldalignment -fix ./...

go-mock:
	go generate ./...

go-cover:
	go test -cover ./...
	
go-swag:
	swag init -g cmd/main.go

# cmd //c tree "%cd%" //F