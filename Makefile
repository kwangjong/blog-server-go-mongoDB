build:
	cd src && docker build -t api-server .

run:
	docker run -d -p 443:443 api-server

test:
	cd src && go test -v
