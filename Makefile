EXEC_NAME_PREFIX=qadfu
BULD_FLAGS=-tags "fts5"
LDLFAGS="-s -w"
CGO_ENABLED=0

.PHONY: all

build: main.main
all: clean build dist

%.main: main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(EXEC_NAME_PREFIX)_linux_amd64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(EXEC_NAME_PREFIX)_linux_arm64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(EXEC_NAME_PREFIX)_windows_amd64.exe
##	GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(EXEC_NAME_PREFIX)_darwin_amd64

dist:
	test -d dist || mkdir -p dist/
	mv $(EXEC_NAME_PREFIX)_* dist/

clean:
	@rm -f $(EXEC_NAME_PREFIX)_*
	@test -d dist && @rm -f dist/$(EXEC_NAME_PREFIX)~
	@test -d dist && @rmdir dist
