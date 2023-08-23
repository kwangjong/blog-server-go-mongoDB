build:
	cd src && go build -o ../server

run:
	sudo ./server

test:
	cd src && go test -v