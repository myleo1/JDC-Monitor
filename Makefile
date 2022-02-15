BINARY=JDC-Monitor
VERSION=1.0.0
DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w"
.PHONY: build deploy

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

deploy:build
	@scp -P 55025 ${BINARY} root@192.168.0.188:~/
	@ssh -p 55025 root@192.168.0.188 "docker-compose -f /root/scripts/docker-app.yml --compatibility stop ${BINARY}"
	@ssh -p 55025 root@192.168.0.188 "mv ${BINARY} /volume2/apps/${BINARY}"
	@ssh -p 55025 root@192.168.0.188 "docker-compose -f /root/scripts/docker-app.yml --compatibility up -d ${BINARY}"
	@echo "[ok] deploy ${BINARY}"

