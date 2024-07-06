FROM alpine:3.13

WORKDIR /root/
COPY profile profile 
COPY env/sample.config env/config
COPY migrations/sql migrations/sql
COPY integrations/email/assets integrations/email/assets

ENTRYPOINT ["/root/profile"]
