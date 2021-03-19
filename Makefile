BINARY=jd-mining-server
VERSION=1.0.0
DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w"
.PHONY: init build build_osx deploy

default:
	@echo ${BINARY}
	@echo ${VERSION}
	@echo ${DATE}

build:
	@GOOS=linux GOARCH=amd64 go build -trimpath -o ${BINARY} ${LDFLAGS}
	@echo "[ok] build ${BINARY}"
	@upx -9 ${BINARY} -o ${BINARY}-upx
	@rm -r ${BINARY}
	@mv ${BINARY}-upx ${BINARY}
