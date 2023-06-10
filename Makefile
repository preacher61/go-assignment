.DEFAULT_GOAL=build

GO_BUILD_DIR=build

.PHONY: build
build:
	mkdir -p $(GO_BUILD_DIR)
	go build -v -o $(GO_BUILD_DIR) ./cmd/...
	