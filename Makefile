db ?= postgres://postgres@localhost/postgres

tools: tools/lint tools/swagger
	
tools/swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

tools/lint: tools/lint/go

tools/lint/go:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install github.com/segmentio/golines@v0.12.2
	go install github.com/daixiang0/gci@v0.13.4
	go install github.com/securego/gosec/v2/cmd/gosec@v2.20.0

lint: tools/lint
	golangci-lint run

lint/fix: tools/swagger
	swag fmt
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default -s "prefix(goal-tracker/api)" .

build: 
	go build -o=./bin/api ./cmd/api
	make swag

run/api:
	go run ./cmd/api

test:
	go test ./cmd/api

test/v:
	go test ./cmd/api -v
	
test/race:
	go test ./cmd/api -race -v

test/pprof:
	go test ./cmd/api -cpuprofile cpu.prof -memprofile mem.prof -bench ./cmd/api

test/cov/report:
	go test ./cmd/api -coverpkg=./cmd/api,./internal/... -covermode=set -coverprofile=coverage.out

test/cov: test/cov/report
	go tool cover -html=coverage.out -o=coverage.html
	make test/cov/open

test/cov/open:
	open ./coverage.html

swag: tools/swagger
	swag init --ot json --parseDependency -g cmd/api/main.go 
	cp -R ./docs/* ./../web/docs 