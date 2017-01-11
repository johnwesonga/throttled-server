# throttled-server
Implements a throttled web server
# Instructions
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server  
docker run --rm -p 3000:3000 jwesonga/throttled-server
