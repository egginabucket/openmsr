ver?=dev
bindir?=./bin
GOARCH?=amd64
GOCMD=./cmd/openmsr/main.go
GOBUILD=go build -ldflags "-w -s"
OUT=$(bindir)/openmsr-$(ver)-$(GOARCH)

clean:
	-rm -rf $(bindir)

run:
	go run $(GOCMD)

dev:
	go build -o $(OUT)-dev $(GOCMD)

linux:
	CGO_ENABLED=1 GOOS=linux
	CC=gcc CXX=g++
ifeq ($(GOARCH), 386)
	GOGCCFLAGS="-m32"
endif
	$(GOBUILD) -o $(OUT) $(GOCMD)

windows:
	CGO_ENABLED=1 GOOS=windows
ifeq ($(GOARCH), 386) 
	CC=i686-w64-mingw32-gcc
	CXX=i686-w64-mingw32-g++
else
	CC=x86_64-w64-mingw32-gcc
	CXX=x86_64-w64-mingw32-g++
endif
	$(GOBUILD) -o $(OUT).exe $(GOCMD)

macos:
	CGO_ENABLED=1 GOOS=darwin
ifeq ($(arch), 386)
	CC=o32-gcc CXX=o32-g++
else
	CC=o64-gcc CXX=064-g++
endif
	$(GOBUILD) -o $(OUT)-macos $(GOCMD)

all:
	make linux windows macos GOARCH=amd64
# 386 doesn't work for anything yet
	-make linux windows macos GOARCH=386