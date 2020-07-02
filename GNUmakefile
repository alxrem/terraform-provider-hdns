BINARY := terraform-provider-hdns
SOURCES := $(wildcard *.go) $(wildcard hdns/*.go)

default: test $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

.PHONY: test
test: $(SOURCES)
	go test ./...

.PHONY: testacc
testacc: $(SOURCES)
	TF_ACC=1 go test ./... -v -timeout 30m

.PHONY: sweep
sweep:
	TF_ACC=1 go test ./... -v -sweep all
