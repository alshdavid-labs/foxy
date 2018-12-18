default: build

PHONY: build
build-windows:
<<<<<<< HEAD
	cd src && go build -o ../bin/foxy cmd/main.go

PHONY: build-windows
build-windows:
	cd src && go build -o ../bin/foxy.exe cmd/main.go
=======
	cd src && go build -o ../foxy cmd/main.go

PHONY: build-windows
build-windows:
	cd src && go build -o ../foxy.exe cmd/main.go
>>>>>>> 0ea382569c4fc06d5c0b7dbcdbb9846634909293
