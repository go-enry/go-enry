# Path to the upstream linguist repo for data syncing
LINGUIST_PATH = .linguist

# Directories
RESOURCES_DIR=./.shared
PYTHON_ENRY_DIR=./python/enry

# shared objects names
LINUX_LIB = libenry.so
DARWIN_LIB = libenry.dylib
HEADER_FILE=libenry.h

LINUX_DIR=$(RESOURCES_DIR)/linux-x86-64
LINUX_SHARED_LIB=$(LINUX_DIR)/${LINUX_LIB}
DARWIN_DIR=$(RESOURCES_DIR)/darwin
DARWIN_SHARED_LIB=$(DARWIN_DIR)/${DARWIN_LIB}
STATIC_LIB=$(RESOURCES_DIR)/libenry.a
NATIVE_LIB=./shared/enry.go

all: shared static

# --- Linguist & Code Generation ---

$(LINGUIST_PATH):
	git clone https://github.com/github/linguist.git $@

clean-linguist:
	rm -rf $(LINGUIST_PATH)

code-generate: $(LINGUIST_PATH)
	mkdir -p data && \
	go run internal/code-generator/main.go
	ENRY_TEST_REPO="$${PWD}/.linguist" go test  -v \
		-run Test_GeneratorTestSuite \
		./internal/code-generator/generator \
		-testify.m TestUpdateGeneratorTestSuiteGold \
		-update_gold

# --- Shared Library Targets ---

# Master shared target - detects host OS
shared:
	@if [ "$$(uname)" = "Darwin" ]; then \
		$(MAKE) darwin-shared; \
	else \
		$(MAKE) linux-shared; \
	fi


# These act as aliases for the file-based targets below
linux-shared: $(LINUX_SHARED_LIB)
darwin-shared: $(DARWIN_SHARED_LIB)

$(LINUX_SHARED_LIB):
	mkdir -p $(LINUX_DIR)
	# Build shared object natively
	CGO_ENABLED=1 GOOS=linux go build -buildmode=c-shared -o $(LINUX_SHARED_LIB) $(NATIVE_LIB)
	mv $(LINUX_DIR)/$(HEADER_FILE) $(RESOURCES_DIR)/$(HEADER_FILE)
	# Copy for Python ABI discovery
	cp $(LINUX_SHARED_LIB) $(PYTHON_ENRY_DIR)
	cp $(RESOURCES_DIR)/$(HEADER_FILE) $(PYTHON_ENRY_DIR)
	@echo "$(LINUX_LIB) built and copied to $(PYTHON_ENRY_DIR)"

$(DARWIN_SHARED_LIB):
	mkdir -p $(DARWIN_DIR)
	# Build shared object. We remove o64-clang to support native M1/M2/M3 builds.
	CGO_ENABLED=1 GOOS=darwin go build -buildmode=c-shared -o $(DARWIN_SHARED_LIB) $(NATIVE_LIB)
	mv $(DARWIN_DIR)/$(HEADER_FILE) $(RESOURCES_DIR)/$(HEADER_FILE)
	# Copy for Python ABI discovery
	cp $(DARWIN_SHARED_LIB) $(PYTHON_ENRY_DIR)
	cp $(RESOURCES_DIR)/$(HEADER_FILE) $(PYTHON_ENRY_DIR)
	@echo "$(DARWIN_LIB) built and copied to $(PYTHON_ENRY_DIR)"

## --- Static Library ---

static:
	mkdir -p $(RESOURCES_DIR)
	CGO_ENABLED=1 go build -buildmode=c-archive -o $(STATIC_LIB) $(NATIVE_LIB)

# --- Benchmarks ---

benchmarks: $(LINGUIST_PATH)
	go test -run=NONE -bench=. && \
	benchmarks/linguist-total.rb

benchmarks-samples: $(LINGUIST_PATH)
	go test -run=NONE -bench=. -benchtime=5us && \
	benchmarks/linguist-samples.rb

benchmarks-slow: $(LINGUIST_PATH)
	mkdir -p benchmarks/output && \
	go test -run=NONE -bench=. -slow -benchtime=100ms -timeout=100h > benchmarks/output/enry_samples.bench && \
	benchmarks/linguist-samples.rb 5 > benchmarks/output/linguist_samples.bench


# --- Cleanup ---

clean-shared:
	rm -rf $(RESOURCES_DIR)
	rm -f $(PYTHON_ENRY_DIR)/*.so $(PYTHON_ENRY_DIR)/*.dylib

clean: clean-linguist clean-shared

.PHONY: all shared linux-shared darwin-shared static clean code-generate benchmarks benchmarks-samples benchmarks-slow
