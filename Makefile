.PHONY:
	start stop build restart

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/num2word ./num2word.go

start:
	docker-compose up

stop:
	docker-compose down

restart:
	docker-compose down
	docker-compose up --build
