# Path to the upstream linguist repo for data syncing
LINGUIST_PATH = .linguist

# Directories
RESOURCES_DIR=./.shared
PYTHON_ENRY_DIR=./python/enry

# shared objects names
LINUX_LIB = libenry.so
DARWIN_LIB = libenry.dylib
HEADER_FILE=libenry.h

STATIC_LIB=$(RESOURCES_DIR)/libenry.a
NATIVE_LIB=./shared/enry.go

# --- Architecture Detection ---
# Use GOARCH if set (by CI), otherwise fallback to native go env
TARGET_ARCH ?= $(shell go env GOARCH)
# Use GOOS if set (by CI), otherwise fallback to native go env
TARGET_OS ?= $(shell go env GOOS)

# Map Go OS to library extensions
ifeq ($(TARGET_OS),darwin)
    LIB_EXT = dylib
else
    LIB_EXT = so
endif

# Define the specific output path based on OS and ARCH
# This prevents different builds from overwriting each other in CI
BUILD_DIR = $(RESOURCES_DIR)/$(TARGET_OS)-$(TARGET_ARCH)
SHARED_LIB = $(BUILD_DIR)/libenry.$(LIB_EXT)

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
shared: $(SHARED_LIB)

$(SHARED_LIB):
	mkdir -p $(BUILD_DIR)
	# CGO_ENABLED=1 is required for buildmode=c-shared
	# We let Go handle GOOS and GOARCH from the environment variables
	CGO_ENABLED=1 go build -v -buildmode=c-shared -o $(SHARED_LIB) $(NATIVE_LIB)
	
	# Move header to a central location
	mv $(BUILD_DIR)/$(HEADER_FILE) $(RESOURCES_DIR)/$(HEADER_FILE)
	
	# Copy to Python directory for packaging
	mkdir -p $(PYTHON_ENRY_DIR)
	cp $(SHARED_LIB) $(PYTHON_ENRY_DIR)/
	cp $(RESOURCES_DIR)/$(HEADER_FILE) $(PYTHON_ENRY_DIR)/
	@echo "Successfully built $(SHARED_LIB) and copied to $(PYTHON_ENRY_DIR)"

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
	rm -f $(PYTHON_ENRY_DIR)/*.so $(PYTHON_ENRY_DIR)/*.dylib $(PYTHON_ENRY_DIR)/*.h

clean: clean-linguist clean-shared

.PHONY: all shared static clean code-generate benchmarks benchmarks-samples benchmarks-slow
