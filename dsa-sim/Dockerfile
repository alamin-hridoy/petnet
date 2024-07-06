FROM alpine:3.13

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY dsa-sim dsa-sim
COPY env/sample.config env/config

ENTRYPOINT ["/root/dsa-sim"]
