EXEC_NAME_PREFIX=qadfu
BULD_FLAGS=-tags "fts5"

.PHONY: all

build: main.main
all: clean build dist

%.main: main.go
	GOOS=linux GOARCH=amd64 go build $(BULD_FLAGS) -o $(EXEC_NAME_PREFIX)_linux_amd64
	GOOS=linux GOARCH=arm64 go build $(BULD_FLAGS) -o $(EXEC_NAME_PREFIX)_linux_arm64
	GOOS=windows GOARCH=amd64 go build $(BULD_FLAGS) -o $(EXEC_NAME_PREFIX)_windows_amd64.exe
##	GOOS=darwin GOARCH=amd64 go build $(BULD_FLAGS) -o $(EXEC_NAME_PREFIX)_darwin_amd64

dist:
	test -d dist || mkdir -p dist/
	mv $(EXEC_NAME_PREFIX)_* dist/

clean:
	@rm -f $(EXEC_NAME_PREFIX)_*
	@test -d dist && @rm -f dist/$(EXEC_NAME_PREFIX)~
	@test -d dist && @rmdir dist
