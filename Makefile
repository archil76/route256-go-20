include make/lint.mk
include make/build.mk

lint: cart-lint loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

up:
	docker-compose up --build -d

run-all:
	docker-compose up --build -d
