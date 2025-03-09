# Order Management
Example order management API with simple web interface

## Run
`go run main.go`

Server will be available at http://localhost:8080

Order data is generated at `data/`

## Build
Builds executable

`go build -o output/`

## Test 
`go list -f '{{.Dir}}/...' -m | xargs go test`