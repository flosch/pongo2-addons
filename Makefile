# Include go binaries into path
export PATH := $(GOPATH)/bin:$(PWD)/bin:$(PATH)

TEST_SOURCE_PATH=TEST_SOURCE_PATH=$(PWD)

install: mod ## Run installing
	@echo "Environment installed"

test: ## Run test with covering
	$(TEST_SOURCE_PATH) go test -coverprofile=$(PWD)/coverage.out .
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

#############################
#############################
#############################

clean-cache: ## Clean golang cache
	@echo "clean-cache started..."
	go clean -cache
	go clean -testcache
	@echo "clean-cache complete!"

clean-vendor: ## Remove vendor folder
	@echo "clean-vendor started..."
	rm -fr ./vendor
	@echo "clean-vendor complete!"

mod: ## Download all dependencies
	@echo "======================================================================"
	@echo "Run MOD...."
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod tidy
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod vendor
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod download
	@echo "======================================================================"

clean-full: clean-vendor
	@echo "Run clean"
	go clean -i -r -x -cache -testcache -modcache
