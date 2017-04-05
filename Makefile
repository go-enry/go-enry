LINGUIST_PATH = .linguist

$(LINGUIST_PATH):
	git clone https://github.com/github/linguist.git $@

test: $(LINGUIST_PATH)
	go test -v ./...

code-generate: $(LINGUIST_PATH)
	go run internal/code-generator/main.go

clean:
	rm -rf $(LINGUIST_PATH)