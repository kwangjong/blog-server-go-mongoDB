FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum vendor/ ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p /usr/local/bin
RUN go build -v -o /usr/local/bin/app .

CMD ["app"]
