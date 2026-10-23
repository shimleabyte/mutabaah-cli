BINARY ?= mutabaah

.PHONY: build test vet fmt-check check clean

build:
	go build -o $(BINARY) .

test:
	go test -race ./...

vet:
	go vet ./...

fmt-check:
	@files="$$(gofmt -l .)"; \
	if [ -n "$$files" ]; then \
		printf 'These Go files need gofmt:\n%s\n' "$$files"; \
		exit 1; \
	fi

check: fmt-check test vet build

clean:
	rm -f "$(BINARY)" "$(BINARY).exe"
