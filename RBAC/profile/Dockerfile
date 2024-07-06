FROM alpine:3.13

WORKDIR /root/
ADD profile profile 
ADD env/sample.config env/config
ADD migrations/sql migrations/sql

ENTRYPOINT ["/root/profile"]
