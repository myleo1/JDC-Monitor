BINARY=JDC-Monitor
BUILD_DIR=build
VERSION=v1.4
DATE=`date +%FT%T%z`
GO_VERSION=`go version`
LDFLAGS="-s -w -X main.version=${VERSION} -X 'main.date=${DATE}' -X 'main.goVersion=${GO_VERSION}'"
.PHONY: build
OS_ARCH=darwin:amd64 darwin:arm64 linux:386 linux:amd64 linux:arm linux:arm64 linux:mips:softfloat linux:mipsle:softfloat linux:mips64 linux:mips64le linux:riscv64 freebsd:386 freebsd:amd64 windows:386 windows:amd64 windows:arm64

default:
	@echo ${BINARY}
	@echo ${VERSION}
	@echo ${DATE}
	@echo ${GO_VERSION}

jdc: clean
	@CGO_ENABLED=0 go build -trimpath -ldflags $(LDFLAGS) -o ./build/bin/${BINARY} main.go

build: clean app

clean:
	@rm -rf ${BUILD_DIR}

app:
	@$(foreach n, $(OS_ARCH),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}-$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags $(LDFLAGS) -o ./build/${BINARY}-$${target_suffix} main.go;\
		echo "Build $${os}-$${arch} done";\
	)
	@mv ./build/${BINARY}-windows-386 ./build/${BINARY}-windows-386.exe
	@mv ./build/${BINARY}-windows-amd64 ./build/${BINARY}-windows-amd64.exe
	@mv ./build/${BINARY}-windows-arm64 ./build/${BINARY}-windows-arm64.exe