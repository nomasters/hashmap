# This make file steamlines the embeded webapp build process

WASM_SOURCE=wasm/main.go
WASM_DEST=src/hashmap.wasm
BUILD_DIR=./public

all: build-wasm build-webapp statik

build-wasm:
	GOOS=js GOARCH=wasm go build -o ${WASM_DEST} ${WASM_SOURCE}

build-webapp: build-wasm
	npm install
	npm run build
	cp ${WASM_DEST} ${BUILD_DIR}/.

statik: build-webapp
	statik -src=${BUILD_DIR}