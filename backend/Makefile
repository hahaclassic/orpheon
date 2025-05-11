ifneq (,$(wildcard ./.env))
    include .env
    export
endif

test:
	go test ./...

# Old approach
# go test -covermode=atomic -coverprofile=./coverage/coverage.out ./internal/... >> /dev/null
# go tool cover -func=./coverage/coverage.out | grep "total:"
cov:
	mkdir -p ./coverage
	./scripts/cov.sh

cov-v: 
	mkdir -p ./coverage
	./scripts/cov.sh -v

integration:
	go test -v ./integration-tests/

clear:
	rm -rf ./coverage

mock:
	mockery --all --recursive --output=./mocks --with-expecter

lint: 
	golangci-lint run
	
migrate-up:
	goose -dir db/migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL_MODE}" up

migrate-down:
	goose -dir db/migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL_MODE}" down

