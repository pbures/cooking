SWAGGGER=./src/swagger/swagger.yml

default: generate

initialize:
	go mod init cooking.buresovi.net

validate:$(SWAGGGER) 
	swagger validate $(SWAGGGER)

generate:
	mkdir -p src/gen-server
	swagger generate server -A Cooking -f $(SWAGGGER) -t src/gen-server --exclude-main

clean:
	rm -f ./main
	rm -rf  ./src/gen-server
 
build: generate
	go build src/main.go

run-server: generate build run-postgres
	./main

run-postgres:
	podman run --name cooking-db -p 5432:5432 -e POSTGRES_PASSWORD=welcome1 -d postgres

stop-postgres:
	podman stop cooking-db
	podman rm cooking-db