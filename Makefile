.PHONY: all test release snapshot

all: test

test:
	go test ./...

clean:
	rm -rf dist

golangci-lint:
	golangci-lint run ./...

## --------------------------------------
## Release
## --------------------------------------

.PHONY: release
release:
	goreleaser release

# used when need to validate the goreleaser
.PHONY: snapshot
snapshot:
	goreleaser release --skip-sign --skip-publish --snapshot --rm-dist
