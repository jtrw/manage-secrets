# Manager secrets

## Use
`go run backend/app/main.go run` - Run web server
`go run backend/app/main.go kv set fkey value` - Set Value
`go run backend/app/main.go kv get fkey` - Get Value


## Use
1. boltdb

### Dev

## Add local package

`cd pkg/utils` => `go mod init`

in main `go.mod` file add our module
```
require (
    ...
	utils v1.0.0 // indirect
)

replace utils v1.0.0 => ./pkg/utils
```

run command `go mod tidy` and get package in main project `go get utils@v1.0.0 `