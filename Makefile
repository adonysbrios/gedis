BINARY=gedis-server
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

.PHONY: all build test lint bench fuzz clean run race

all: clean build test lint

build:
	$(GOBUILD) -o $(BINARY) .

test:
	$(GOTEST) -count=1 -race ./...

bench:
	$(GOTEST) -bench=. -benchmem ./...

lint:
	$(GOLINT) run ./...

fuzz:
	$(GOTEST) -fuzz=FuzzParseArgs -fuzztime=10s ./protocol
	$(GOTEST) -fuzz=FuzzParseRESP -fuzztime=10s ./protocol

race:
	$(GOTEST) -race -count=1 ./...

vet:
	$(GOVET) ./...

run:
	$(GOBUILD) -o $(BINARY) .
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -f gedis.aof gedis.aof.tmp
