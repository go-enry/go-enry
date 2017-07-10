COVERAGE_REPORT := coverage.txt
COVERAGE_PROFILE := profile.out
COVERAGE_MODE := atomic

LINGUIST_PATH = .linguist

# build CLI
VERSION := $(shell git describe --tags --abbrev=0)
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS = -s -X main.Version=$(VERSION) -X main.GitHash=$(COMMIT)

$(LINGUIST_PATH):
	git clone https://github.com/github/linguist.git $@

test: $(LINGUIST_PATH)
	go test -v ./...

test-coverage: $(LINGUIST_PATH)
	@echo "mode: $(COVERAGE_MODE)" > $(COVERAGE_REPORT); \
		for dir in `find . -name "*.go" | grep -o '.*/' | sort -u | grep -v './fixtures/' | grep -v './.linguist/'`; do \
			go test $$dir -coverprofile=$(COVERAGE_PROFILE) -covermode=$(COVERAGE_MODE); \
			if [ $$? != 0 ]; then \
				exit 2; \
			fi; \
			if [ -f $(COVERAGE_PROFILE) ]; then \
				tail -n +2 $(COVERAGE_PROFILE) >> $(COVERAGE_REPORT); \
				rm $(COVERAGE_PROFILE); \
			fi; \
	done;

code-generate: $(LINGUIST_PATH)
	mkdir -p data
	go run internal/code-generator/main.go

benchmarks: $(LINGUIST_PATH)
	go test -run=NONE -bench=. && benchmarks/linguist-total.sh

benchmarks-samples: $(LINGUIST_PATH)
	go test -run=NONE -bench=. -benchtime=5us && benchmarks/linguist-samples.rb

benchmarks-slow: $(LINGUST_PATH)
	mkdir -p benchmarks/output && go test -run=NONE -bench=. -slow -benchtime=100ms -timeout=100h >benchmarks/output/enry_samples.bench && \
	benchmarks/linguist-samples.rb 5 >benchmarks/output/linguist_samples.bench

clean:
	rm -rf $(LINGUIST_PATH)

build-cli:
	go build -o enry -ldflags "$(LDFLAGS)" cli/enry/main.go
