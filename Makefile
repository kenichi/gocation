all: clean build

db_url: all
	$(eval DB_URL := $(shell ./gocation -db-url))

migrate: db_url
	migrate -url $(DB_URL) -path ./migrations up

reset: db_url
	migrate -url $(DB_URL) -path ./migrations reset
	
clean:
	go clean

build:
	go build
