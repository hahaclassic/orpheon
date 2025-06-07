#!/bin/bash

go test -covermode=atomic -coverprofile=./coverage/coverage.out ./internal/... ./cmd/... ./pkg/... > coverage/result.log

# go test -covermode=atomic -coverprofile=./coverage/coverage.out ./internal/domain/services/... ./internal/adapters/... > coverage/result.log

if [ $? -eq 0 ]; then
    if [[ "$1" == "-v" ]]; then 
        cat coverage/result.log
    fi
    coverage=$(go tool cover -func=./coverage/coverage.out | grep "total:" | awk '{print $3}')
    echo "status: OK. coverage: $coverage"
else
    echo "status: FAIL"
    cat coverage/result.log
fi
