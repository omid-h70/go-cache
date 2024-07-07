test-status:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Test all test functions that exist in this package
test:
	go test -v ./...

test-cover:
	go test -cover .