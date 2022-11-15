FROM alpine:latest
LABEL maintainer="leo <leo@leosgo.com>" \
	version="v1.0.0" \
	description="JDC-Monitor"
WORKDIR /root
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
	&& apk update \
	&& apk add --no-cache tzdata
ADD localtime /etc/localtime
ADD JDC-Monitor /root/JDC-Monitor
CMD ["/root/JDC-Monitor"]