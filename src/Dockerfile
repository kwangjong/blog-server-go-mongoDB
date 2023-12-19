FROM golang:1.21

WORKDIR /usr/src/app

COPY . .
RUN cat 107106.xyz-ssl-bundle/domain.cert.pem 107106.xyz-ssl-bundle/intermediate.cert.pem > 107106.xyz-ssl-bundle/combined.cert.pem

RUN mkdir -p /usr/local/bin
RUN go build -v -o /usr/local/bin/app .

CMD ["app"]
