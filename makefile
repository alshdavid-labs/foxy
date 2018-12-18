default: build

PHONY: build
build-windows:
	cd src && go build -o ../bin/foxy cmd/main.go

PHONY: build-windows
build-windows:
	cd src && go build -o ../bin/foxy.exe cmd/main.go