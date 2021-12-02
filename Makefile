CGO_ENABLED=0
GO_BUILD_FLAGS=-ldflags "-extldflags '-static'"

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o bin/app ./cmd/app

.PHONY: clean
clean:
	rm -rf ./bin/app

.PHONY: format
format:
	go fmt $(go list ./... | grep -v /vendor/)

.PHONY: test
test:
	go vet $(go list ./... | grep -v /vendor/)
	go test -race $(go list ./... | grep -v /vendor/)
