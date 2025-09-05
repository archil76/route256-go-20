BINDIR=${CURDIR}/bin
PACKAGE=route256/cart

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build cart"
	docker-compose up --build

run-all: bindir
	echo "build cart"
	docker-compose up --build

.install-minimock:
	go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5

gen-mock:
	minimock -i ./internal/domain/service.CartsRepository -o ./internal/domain/service/mock/cart_repository_mock.go
	minimock -i ./internal/domain/service.ProductService -o ./internal/domain/service/mock/product_service_mock.go

test-all:
	go test -cover -tags e2e_test ./internal/app/server ./internal/domain/repository/inmemoryrepository

bench-all:
	go test ./internal/domain/repository/inmemoryrepository -bench=.

test-serv: run-all
	export CONFIG_FILE=./cart/configs/values_local.yaml
	rm ./testы/e2e/allure-results/*
	go test ./tests/e2e

allure-serve:
	allure serve ./tests/e2e/allure-results

.install-linters:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/uudashr/gocognit/cmd/gocognit@latest

lint:
	gocyclo -top 10 -avg -ignore "mock" .
	gocognit -top 10 -avg -ignore "mock" .
	golangci-lint run


