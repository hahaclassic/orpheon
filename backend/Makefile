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

clear:
	rm -rf ./coverage

gen-mocks:
	mockery --all --recursive --output=./mocks --with-expecter

lint: 
	golangci-lint run
	