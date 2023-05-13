VERSION=v0.0.1

bin: bin/t7s_darwin bin/t7s_linux bin/t7s_windows

bin/t7s_darwin:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/t7s_darwin cmd/t7s/*.go
	openssl sha512 bin/t7s_darwin > bin/t7s_darwin.sha512

bin/t7s_linux:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/t7s_linux cmd/t7s/*.go
	openssl sha512 bin/t7s_linux > bin/t7s_linux.sha512

bin/t7s_windows:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/t7s_windows cmd/t7s/*.go
	openssl sha512 bin/t7s_windows > bin/t7s_windows.sha512
