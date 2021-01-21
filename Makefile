test:
	go test -v ./employee/ ./project/
build:
	docker-compose -f docker-compose.yml build
up:
	@make build
	docker-compose -f docker-compose.yml up -d
logs:
	docker-compose logs -f
rm:
	docker-compose rm  -sfv
start:
	docker-compose start
stop:
	docker-compose stop
rest:
	@make rm
	@make up
	@make logs

bash_db:
	docker-compose run database bash
bash_go:
	docker-compose run go sh