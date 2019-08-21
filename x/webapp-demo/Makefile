# This make file steamlines the embeded webapp build process

WASM_SOURCE=wasm/main.go
WASM_OUT=assets/hashmap.wasm
BUILD_DIR=./public

all: build-wasm build-webapp statik

build-wasm:
	GOOS=js GOARCH=wasm go build -o ${BUILD_DIR}/${WASM_OUT} ${WASM_SOURCE}

build-webapp: build-wasm
	npm install
	npm run build

statik: build-webapp
	statik -src=${BUILD_DIR}