go-tests := "go test ./..."
go-lint := "GOFLAGS=-buildvcs=false golangci-lint run"

watch-dev:
   just --justfile {{justfile()}} init-git-hooks
   just --justfile {{justfile()}} tidy
   watchexec -r -e go 'go run main.go'

watch-tests:
   watchexec -e go '{{go-tests}}'

watch-lint:
   watchexec -e go '{{go-lint}}'

tidy:
   go mod tidy

init-git-hooks:
    cp pre-commit.sh .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit

pre-commit:
    @echo "Running pre-commit hook..."
    {{go-lint}}
    {{go-tests}}
