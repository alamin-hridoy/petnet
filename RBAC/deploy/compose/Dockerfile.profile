FROM golang:1.16-alpine

WORKDIR /src

ENV VERSION=localdev

COPY . .

WORKDIR profile 

RUN go build -mod=vendor -ldflags="-w -s -X main.version=$VERSION" -o profile && \
	if [ ! -f env/config ]; then cp env/sample.config env/config ; fi

CMD ["sh", "-c", "./profile"]
