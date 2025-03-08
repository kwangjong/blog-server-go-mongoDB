FROM golang:1.24

ENV GOFLAGS="-mod=vendor"

WORKDIR /app

COPY . .
RUN mkdir -p /usr/local/bin
RUN go build -v -o /usr/local/bin/app .

CMD ["app"]
