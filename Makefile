SWAGGGER=./src/swagger/swagger.yml

default: generate

initialize:
	go mod init cooking.buresovi.net

validate:$(SWAGGGER) 
	swagger validate $(SWAGGGER)

generate:
	mkdir -p src/server
	swagger generate server -A Cooking -f $(SWAGGGER) -t src/server

clean:
	rm -f ./main
	cd src/server && rm -rf cmd models
	cd src/server/restapi && rm -rf operations doc.go server.go embedded_spec.go
 
build: generate
	go build src/server/cmd/cooking-server/main.go

run-server: generate build run-postgres
	./main

run-postgres:
	podman run --name cooking-db -p 5432:5432 -e POSTGRES_PASSWORD=welcome1 -d postgres

stop-postgres:
	podman stop cooking-db
	podman rm cooking-db