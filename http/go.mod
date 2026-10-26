module github.com/ifnotnil/x/http

go 1.26.0

require (
	github.com/andybalholm/brotli v1.2.7
	github.com/klauspost/compress v1.20.3
	golang.org/x/text v0.43.0
)

// Test dependencies. They will not be pushed downstream as indirect ones.
require github.com/stretchr/testify v1.12.1

require (
	github.com/stretchr/objx v0.8.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)
