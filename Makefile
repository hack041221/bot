CGO_ENABLED=0
GO_BUILD_FLAGS=-ldflags "-extldflags '-static'"

.PHONY: docker-build-bot
docker-build-bot:
	docker build -t hack-bot -f bot.dockerfile .

.PHONY: docker-build-downloader
docker-build-downloader:
	docker build -t hack-downloader -f downloader.dockerfile .

.PHONY: build-bot
build-bot:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o bin/app ./cmd/app

.PHONY: build-downloader
build-downloader:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o bin/download ./cmd/download

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
