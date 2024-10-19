
.PHONY: test
test:
	find . -name go.mod -execdir go test -race -v ./... \;

.PHONY: cov
cov:
	find . -name go.mod -execdir go test -race -covermode=atomic -coverprofile=coverage.out ./... \;

.PHONY: tidy
tidy:
	find . -name go.mod -execdir go mod tidy \;
