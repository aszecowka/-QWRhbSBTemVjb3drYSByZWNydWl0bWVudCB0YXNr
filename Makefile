APP_NAME = "bulk-fetcher"

.PHONY: validate
validate:
	go mod tidy -v
	go fmt ./...
	go vet ./...
	go install github.com/kisielk/errcheck
	errcheck -verbose ./...

.PHONY: tests
tests: validate
	go test ./...

.PHONY: build
build: tests
	go build -o ${APP_NAME} ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: run-with-deps
run-with-deps:
	docker-compose build server
	docker-compose up

.PHONY: run-int-test
run-int-test:
	docker-compose -f docker-compose-int-test.yml build server integration-test
	docker-compose -f docker-compose-int-test.yml up integration-test

.PHONY: build-docker
build-docker:
	docker build -t ${APP_NAME} .

.PHONY: build-docker-int
build-docker-int:
	docker build -t ${APP_NAME}-int-test -f Dockerfile-integration-test .
