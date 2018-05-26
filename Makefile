INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
TEST_MODEL ?= github.com/sonod/pisera/model
TEST_OTHERS ?= $(shell go list ./... | grep -v vendor | grep -v api | grep -v model)
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)
VERBOSE = $(shell test -z $$DEBUG && test -z $$VERBOSE || echo '-v')

build: ## Build as linux binary
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building$(RESET)"
	./misc/build

test: testmodel testothers testrace ## Run test

testmodel: ## Run model test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing Model$(RESET)"
	go test $(VERBOSE) $(TEST_MODEL) $(TEST_OPTIONS)
	go test $(VERBOSE) $(TEST_MODEL)/model_test $(TEST_OPTIONS)

testothers: ## Run other tests
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing Others$(RESET)"
	go test $(VERBOSE) $(TEST_OTHERS) $(TEST_OPTIONS)

testrace: ## Run race test
	go test $(VERBOSE) -race $(TEST_MODEL)
	go test $(VERBOSE) -race $(TEST_API)
	go test $(VERBOSE) -race $(TEST_OTHERS)
