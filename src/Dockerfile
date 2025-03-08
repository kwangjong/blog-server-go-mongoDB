FROM golang:1.24

ENV GOFLAGS="-mod=vendor"

WORKDIR /app

COPY . .
RUN ls -l vendor/
RUN go clean -modcache
RUN go build -mod=vendor -v -o app .

CMD ["app"]
