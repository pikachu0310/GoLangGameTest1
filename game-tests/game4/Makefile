.PHONY: build-wasm build-windows copy-wasm-exec

# APIKEYをハードコーディングにしてからwasmに変換する
build-wasm:
	GOOS=js GOARCH=wasm go build -o client/game.wasm game4 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o game.exe .

copy-wasm-exec:
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js client
