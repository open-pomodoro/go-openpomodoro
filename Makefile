default: test lint

test:
	go test ./...

lint:
	@which -s gometalinter || (go get github.com/alecthomas/gometalinter && gometalinter --install)
	gometalinter

.PHONY: test lint
