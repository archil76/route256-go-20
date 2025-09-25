include make/lint.mk
include make/build.mk

lint: cart-lint loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

up:
	docker-compose up --build -d
	cd loms && make db-migrate && cd ..

down:
	cd loms && make db-migrate-down
	cd .. && docker-compose down

run-all:
	docker-compose up --build -d

