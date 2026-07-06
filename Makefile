GO ?= go
IMAGE ?= devbox:local

.PHONY: all test build tidy fmt vet check-loc docker-build check clean

all: check

test:
	$(GO) test ./...

build:
	mkdir -p bin
	$(GO) build -trimpath -ldflags="-s -w" -o bin/devbox ./cmd/devbox

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

check-loc:
	./scripts/check-loc.sh

check: fmt tidy check-loc test vet build

docker-build:
	docker build -t $(IMAGE) .

clean:
	rm -rf bin dist tmp coverage.out
