ARG CHROME_VERSION
ARG GO_VERSION

FROM chromedp/headless-shell:${CHROME_VERSION} as chrome

FROM golang:${GO_VERSION}

COPY --from=chrome /headless-shell/ /headless-shell/
ENV PATH /headless-shell:$PATH
RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
	libnspr4 libnss3 libexpat1 libfontconfig1 libuuid1 \
	ca-certificates wget openssl libnss3-tools \
	&& rm -rf /var/lib/apt/lists/*
