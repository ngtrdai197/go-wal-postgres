.PHONY: build_backend build up down api

build_backend:
	docker-compose build backend

build:
	docker-compose up --build

up:
	docker-compose up

down:
	docker-compose down

api:
	go run main.go public-api