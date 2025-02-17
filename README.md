# run server
```go
go run cmd/server/main.go
```
# run server with params
```go
go run cmd/server/main.go -a 9999
```

# run client
```go
go run cmd/client/main.go
```
# run server with params
```go
go run cmd/client/main.go -a 9999 -p 3 -r 3
```

# run all tests
```go
go test -v ./internal/server/
```