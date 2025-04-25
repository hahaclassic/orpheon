test:
	go test ./...

cov:
	mkdir -p ./coverage
	go test -covermode=atomic -coverprofile=./coverage/coverage.out ./internal/... >> /dev/null
	go tool cover -func=./coverage/coverage.out | grep "total:"

clear:
	rm -rf ./coverage

gen-mocks:
	mockery --all --recursive --output=./mocks --with-expecter