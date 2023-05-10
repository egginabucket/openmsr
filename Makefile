build:
	go build -o ./bin/openmsr ./cmd/openmsr/main.go

compile:
	GOOS=linux GOARCH=amd64 go build -o ./bin/openmsr-amd64 ./cmd/openmsr/main.go
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o ./bin/openmsr-amd64.exe ./cmd/openmsr/main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ./bin/openmsr-macOS-amd64 ./cmd/openmsr/main.go