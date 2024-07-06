FROM alpine:3.13

WORKDIR /root/

RUN apk add --no-cache ca-certificates tzdata \
	&& update-ca-certificates

COPY idp idp
COPY env/sample.config env/config

ENTRYPOINT ["./idp"]
