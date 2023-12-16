# Guestbook-example

#### Using Go version 1.21

## Run tests

```bash
# Install mockery
go install github.com/vektra/mockery/v2@v2.38.0

# Generate mock code
go generate ./...

# Run tests
go test -v ./...
```

## Run code

```bash
go build guestbook-example && ./guestbook-example
```
