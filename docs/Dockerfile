FROM asia.gcr.io/o-production/redoc-builder:v21.06.03 AS builder
WORKDIR /tmp
# api.yaml is built in CI
COPY theme.args api.yaml template.hbs ./
# we only have a handful of args, so no need to worry about redoc-cli being called more than once
RUN cat theme.args | xargs redoc-cli bundle api.yaml -t template.hbs -o redoc.html

FROM alpine:3.13
WORKDIR /root/
COPY assets/ assets/
COPY doc-portal doc-portal
COPY env/sample.config env/config
COPY --from=builder /tmp/redoc.html assets/templates/redoc.html

ENTRYPOINT ["/root/doc-portal"]
