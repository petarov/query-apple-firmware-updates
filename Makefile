.PHONY: all

build: main.main
all: clean build dist

%.main: main.go
	GOOS=linux GOARCH=amd64 go build -o qados_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o qados_linux_arm64
#	GOOS=darwin GOARCH=amd64 go build -o qados_darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o qados_windows_amd64.exe	

dist:
	test -d dist || mkdir -p dist/
	mv qados_* dist/

clean:
	@rm -f qados_*
	@test -d dist && @rm -f dist/qados~
	@test -d dist && @rmdir dist
	