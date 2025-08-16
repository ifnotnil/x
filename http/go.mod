module github.com/ifnotnil/x/http

go 1.24.0

require (
	github.com/andybalholm/brotli v1.2.0
	github.com/klauspost/compress v1.18.0
	golang.org/x/text v0.28.0
)

// Test dependencies. They will not be pushed downstream as indirect ones.
require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
