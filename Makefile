build:
	cd src && go build -o ../server

run:
	./server

test:
	cd src && go test -v
