FROM golang:1.19 as builder
WORKDIR /go/src/JDC-Monitor
COPY . .
RUN make jdc

FROM alpine:latest
WORKDIR /root
RUN apk update \
	&& apk add --no-cache tzdata
COPY localtime /etc/localtime
COPY --from=builder /go/src/JDC-Monitor/build/bin/JDC-Monitor /root/JDC-Monitor
CMD ["/root/JDC-Monitor"]