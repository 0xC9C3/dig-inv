go-tests := "go test -coverprofile=coverage.profile ./..."
go-coverage := "go tool cover -html=coverage.profile -o coverage.html"
go-lint := "GOFLAGS=-buildvcs=false golangci-lint run"

watch-dev-server:
   just --justfile {{justfile()}} init
   watchexec -r -e go 'go run main.go server'

watch-dev-worker:
   just --justfile {{justfile()}} init
   watchexec -r -e go 'go run main.go worker'

watch-tests:
   watchexec -e go '{{go-tests}} && {{go-coverage}}'

watch-lint:
   watchexec -e go '{{go-lint}}'

watch-grpc-buf-generate:
   just --justfile {{justfile()}} init
   watchexec -e proto 'just --justfile {{justfile()}} grpc-buf-generate'

tidy:
   go mod tidy

grpc-buf-generate:
    buf generate

init:
   just --justfile {{justfile()}} init-git-hooks
   just --justfile {{justfile()}} install-tools
   just --justfile {{justfile()}} tidy

install-tools:
    go install tool

init-git-hooks:
    cp pre-commit.sh .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit

pre-commit:
    @echo "Running pre-commit hook..."
    {{go-lint}}
    {{go-tests}}
