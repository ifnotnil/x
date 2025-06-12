MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")

GO_PACKAGES = go list -tags='$(TAGS)' ./...
GO_FOLDERS = go list -tags='$(TAGS)' -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go'

export GO111MODULE := on

.PHONY: mod
mod:
	go mod tidy
	go mod verify

.PHONY: go-gen
go-gen:
	go generate ./...


# https://go.dev/ref/mod#go-get
# -u flag tells go get to upgrade modules
# -t flag tells go get to consider modules needed to build tests of packages named on the command line.
# When -t and -u are used together, go get will update test dependencies as well.
.PHONY: go-deps-upgrade
go-deps-upgrade:
	go get -u -t ./...
	go mod tidy

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

.PHONY: test
test:
	CGO_ENABLED=1 go test -timeout 60s -race -tags='$(TAGS)' -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-n-read
test-n-read: test
	@go tool cover -func coverage.txt

.PHONY: bench
bench: # runs all benchmarks
	CGO_ENABLED=1 go test -benchmem -run=^Benchmark$$ -mod=readonly -count=1 -v -race -bench=. ./...

.PHONY: caches-info
caches-info:
	@echo "cache folders: "
	@du -h -d 0 -a "$$(go env GOMODCACHE)"
	@du -h -d 0 -a "$$(go env GOCACHE)"
