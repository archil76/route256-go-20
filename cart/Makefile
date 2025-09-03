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
	go test -cover ./internal/app/server ./internal/domain/repository/inmemoryrepository




