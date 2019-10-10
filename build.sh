env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o goMongoTest-x86 .
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o goMongoTest-macos .
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o goMongoTest-win32.exe .

docker build -t jerrymannel/gomongotest:1 .
dim tag jerrymannel/gomongotest:1 jerrymannel/gomongotest