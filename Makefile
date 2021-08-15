GOOS ?= $(shell go env GOOS)
SOURCES := $(shell find . -type f  -name '*.go')
LDFLAGS = -ldflags="-s -w" -installsuffix cgo

build: $(SOURCES)
	CGO_ENABLED=0 GOOS=$(GOOS) go build \
	$(LDFLAGS) \
	-o app main.go

clean:
	rm -rf app
