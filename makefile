default: build

PHONY: build
build-windows:
	cd src && go build -o ../foxy cmd/main.go

PHONY: build-windows
build-windows:
	cd src && go build -o ../foxy.exe cmd/main.go