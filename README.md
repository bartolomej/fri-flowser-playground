# Flowser Playground (demo)

This is an example implementation of the [Flowser playground](https://github.com/onflow/developer-grants/issues/260), which will replace the current [Flow playground](https://play.flow.com/).

> Note: This was a prototype project. Development continued on [flow-wasm](https://github.com/onflowser/flow-wasm) and [flowser](https://github.com/onflowser/flowser) repos.

## Features

- Clone project from Github
- Deploy the project contracts
- View and edit project files
- Execute transactions and scripts
- View project logs and blockchain state

<img src="https://github.com/bartolomej/fri-flowser-playground/assets/36109955/a028462e-bf11-4e29-bdbf-a282806d6669" />


## Get started

Install dependencies:
```
cd web && npm i
```

Start backend:

```bash
go run cmd/main.go
```

Start client:

```bash
cd web && npm run dev
```

## Building

### Cross compiling for Windows

Since the dependency [onflow/crypto](https://github.com/onflow/crypto/tree/e9ca850f06dfd0e3f56fe0e3233c1ebb32b2e4d0) depends on native C libraries and uses [cgo](https://go.dev/wiki/cgo) to build that, we must have a C compiler installed locally.

When cross-compiling from MacOS, you must install the MinGW toolchain with:

```bash
brew install mingw-w64
```

Then you can build the program for Windows with:

```bash
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build cmd/main.go
```
