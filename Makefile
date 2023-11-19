run:
	clear
	go run .

build:
	go clean -x
	go clean -cache -x
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -v -ldflags "-w -s -extldflags '-static'" \
	-gcflags="-S -m" -trimpath -mod=readonly -buildmode=pie -a -installsuffix nocgo \
	-o gojo-scraper .
	clear
	./gojo-scraper


.PHONY: run, build