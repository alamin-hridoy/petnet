FROM alpine:3.13

WORKDIR /root/
COPY usermgm usermgm
COPY integrations/email/assets integrations/email/assets
COPY env/sample.config env/config
COPY migrations/sql migrations/sql

ENTRYPOINT ["/root/usermgm"]
