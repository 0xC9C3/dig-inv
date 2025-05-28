watch-dev:
   just --justfile {{justfile()}} tidy
   watchexec -r -e go 'go run main.go'

watch-tests:
   watchexec -e go 'go test .'

tidy:
   go mod tidy
