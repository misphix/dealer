.PHONY: all start stop clean test deploy

all:
	go build

env:
	docker-compose up database -d

env-down:
	docker-compose down

clean:
	rm -f dealer
	rm -f *.out

test:
	go test -v -covermode=count -coverprofile=test.out ./...
	go tool cover -html=test.out

deploy:
	docker-compose build
	docker-compose up