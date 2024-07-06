# syntax = docker/dockerfile:1.3
FROM golang:1.17-alpine

WORKDIR /src

ENV VERSION=localdev

COPY . .

WORKDIR profile

RUN --mount=type=cache,target=/root/.cache/go-build go build --mod=vendor -ldflags="-w -s -X main.version=$VERSION" -o profile && \
	go build --mod=vendor -ldflags="-w -s -X main.version=$VERSION" -o migrate ./migrations && \
	if [ ! -f env/config ]; then cp env/sample.config env/config ; fi

CMD ["sh", "-c", "./migrate up && ./profile"]
