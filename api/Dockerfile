FROM alpine:3.13

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
ADD env/sample.config env/config
ADD migrations/sql migrations/sql

COPY api api

ENTRYPOINT ["sh", "-c", "./api"]
