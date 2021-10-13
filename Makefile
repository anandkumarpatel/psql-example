APP_NAME = psql-example

DB_URL = postgres://postgres:secret@localhost:5432/postgres

default: run

test: dependencies-cleanup dependencies-start
	DB_URL=$(DB_URL) make unit-test
	make dependencies-cleanup

unit-test:
	go test -v -failfast ./...

dependencies-cleanup: 
	docker kill postgres || true
	docker rm postgres || true

dependencies-start: dependencies-cleanup
	docker run --name postgres -p 5432:5432 -v $(shell pwd)/services.sql:/docker-entrypoint-initdb.d/services.sql -e POSTGRES_PASSWORD=secret -d postgres
	until docker exec -it postgres psql  -U postgres -c "select 1"; do sleep 1; done

build:
	go build -o $(APP_NAME) .

run: dependencies-start build
	DB_URL=$(DB_URL) ./$(APP_NAME)
