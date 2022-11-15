FROM alpine:latest
LABEL maintainer="leo <leo@leosgo.com>" \
	version="v1.0.0" \
	description="JDC-Monitor"
WORKDIR /root
ADD JDC-Monitor /root/JDC-Monitor
ADD localtime /etc/localtime
CMD ["/root/JDC-Monitor"]