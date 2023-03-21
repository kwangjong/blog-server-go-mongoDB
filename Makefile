build:
	go build
	go build ./...

run:
	./kwangjong.github.io

test:
	go test ./... -v