GO      ?= go
GOOS    ?= $(shell $(GO) env GOOS)
GOARCH  ?= $(shell $(GO) env GOARCH)
PKG     := ./...
BINDIR  := build/$(GOOS)/$(GOARCH)
BIN     := $(BINDIR)/rt
PGO     := $(GOOS)-$(GOARCH).pgo

.PHONY: all lint test bench profile build clean

all: lint test profile build

build: $(BINDIR)/rt

profile: $(PGO)

$(BINDIR):
	mkdir -p $(BINDIR)

$(BINDIR)/rt: $(BINDIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -pgo $(PGO) -o $(BIN) $(PKG)

$(PGO):
	$(GO) run $(PKG) -width 256 -height 144 -cpuprofile $(PGO) -output /dev/null

clean:
	rm -rf $(BINDIR)

test:
	$(GO) test $(PKG)

lint:
	golangci-lint run --fix

bench:
	$(GO) test -bench=. -run=^$$ $(PKG)
