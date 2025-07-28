.PHONY: run up stop

run: stop up

up:
	docker compose up -d --build

stop:
	docker compose stop
