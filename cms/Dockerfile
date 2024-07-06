FROM alpine:3.13

WORKDIR /root/

RUN apk add --no-cache ca-certificates tzdata \
	&& update-ca-certificates

COPY cms cms
COPY assets assets
COPY env/sample.config env/config

ENTRYPOINT ["sh", "-c", "./cms"]
