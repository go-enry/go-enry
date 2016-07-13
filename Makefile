samples:
	git clone git@github.com:github/linguist.git .linguist\

test: samples
	go test -v ./...