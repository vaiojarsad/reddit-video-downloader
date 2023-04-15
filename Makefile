default: help

.PHONY: help
help: ## Show this help
	@echo
	@echo "Available commands:"
	@echo
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST)

.PHONY: setup
setup: ## Setup system for local development
	go get github.com/spf13/cobra

.PHONY: build
build: ## Just build the program
	go build -o ./bin/reddit-video-downloader.exe ./.development/main.go

.PHONY: clean
clean: ## Clean files
	@-rm ./bin/reddit-video-downloader.exe

.PHONY: tidy
tidy: ## Clean cache and run tidy
	(go clean -modcache && go mod tidy)
