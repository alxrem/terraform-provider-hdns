BINARY := terraform-provider-hdns
SOURCES := $(wildcard *.go) $(wildcard hdns/*.go)

default: test $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

test: $(SOURCES)
	go test ./...

testacc: $(SOURCES)
	TF_ACC=1 go test ./... -v -timeout 30m

sweep:
	TF_ACC=1 go test ./... -v -sweep all

release:
	goreleaser release --rm-dist

snapshot:
	goreleaser release --snapshot --rm-dist --skip-publish

.PHONY: test testsacc sweep release snapshot