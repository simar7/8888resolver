all:
	make clean test fmt lint vet build run

clean:
	rm 8888resolver

test:
	go test -v -cover

fmt:
	go fmt $(go list ./... | grep -v /vendor/)

lint:
	golint $(go list ./... | grep -v /vendor/)

vet:
	go vet -all -shadowstrict $(go list ./... | grep -v /vendor/)

build:
	go build .

run:
	./8888resolver

runprod:
	GIN_MODE=release ./8888resolver