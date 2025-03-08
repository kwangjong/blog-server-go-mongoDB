FROM golang:1.24

WORKDIR /usr/src/app

COPY . .
RUN mkdir -p /usr/local/bin
RUN go build --build-arg GOFLAGS="-mod=vendor" -t -v -o /usr/local/bin/app .

CMD ["app"]
