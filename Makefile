VERSION=v0.0.2

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

robertlestak/t7s:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t robertlestak/t7s:$(VERSION) --push .
	docker buildx imagetools create robertlestak/t7s:$(VERSION) --tag robertlestak/t7s:latest