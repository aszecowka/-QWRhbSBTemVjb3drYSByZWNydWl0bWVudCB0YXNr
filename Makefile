APP_NAME = "bulk-fetcher"

.PHONY: build
build: tests
	go build -o ${APP_NAME} ./cmd/main.go

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

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: run-with-deps
run-with-deps:
	docker-compose up --build

.PHONY: run-int-test
run-int-test:
	docker-compose -f docker-compose-int-test.yml up --build --abort-on-container-exit --exit-code-from integration-test integration-test
	docker-compose down

.PHONY: build-docker
build-docker:
	docker build -t ${APP_NAME} .

.PHONY: build-docker-int
build-docker-int:
	docker build -t ${APP_NAME}-int-test -f Dockerfile-integration-test .
