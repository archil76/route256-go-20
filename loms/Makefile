BINDIR=${CURDIR}/bin
PACKAGE=route256/loms

PROTO_PATH := "./api"

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build loms"

run-migrations:
	echo "run migrations"

gen-mock:
	minimock -i ./internal/app/server.LomsServise -o ./internal/app/server/mock/loms_servise_mock.go

test-all:
	go test -cover ./internal/app/server

bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

protoc-generate: bin-deps
	protoc \
	-I ./api \
	--go_out=./internal/api \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/api \
	--go-grpc_opt=paths=source_relative \
	api/loms.proto && \
	go mod tidy