.PHONY: install-tools format sort-imports run test coverage build-lambda run-all

install-tools:
	go install github.com/daixiang0/gci@latest
	go install mvdan.cc/gofumpt@latest

format:
	gofmt -w .
	gofumpt -l -w .
	go mod tidy

sort-imports:
	@gci write --skip-generated -s standard -s default -s "prefix(github.com/)" -s "prefix(`(head -n 1 ./go.mod | sed 's/^module //')`)" .

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic -v
	go tool cover -func=coverage.out
	@COVERAGE=$$(go tool cover -func=coverage.out | tail -1 | awk '{gsub(/%/,""); print $$3}'); \
	if [ "$$(printf "%.0f" "$$COVERAGE")" -lt 90 ]; then echo "Coverage $$COVERAGE% is below 90%"; exit 1; fi

build-lambda:
	GOOS=linux GOARCH=arm64 go build -o bootstrap .
	zip -q function.zip bootstrap

run-all: format sort-imports test
