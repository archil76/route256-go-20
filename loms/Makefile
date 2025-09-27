BINDIR=${CURDIR}/bin
PACKAGE=route256/loms
CONFIG_FILE=loms/configs/values_local.yaml

MIGRATIONS_FOLDER := ./migrations/
LOCAL_DB_DSN := postgresql://loms-user:loms-password@127.0.0.1:5433/loms_db?sslmode=disable

PROTO_PATH := "./api"

bindir:
	mkdir -p ${BINDIR}

run:
	export CONFIG_FILE=./configs/values_local.yaml &&\
	go run ./cmd/server/main.go

run-all: bindir
	echo "build cart"
	docker-compose up --build

build: bindir bin-deps vendor-proto protoc-generate
	echo "build loms"

run-migrations:
	echo "run migrations"

.install-minimock:
	go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5

gen-mock:
	minimock -i ./internal/app/server.LomsServise -o ./internal/app/server/mock/loms_servise_mock.go

test-all:
	go test -cover ./internal/app/server ./internal/domain/service

bin-deps:
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
	GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
    GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1

vendor-rm:
	rm -rf vendor-proto

vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
	https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
	cd vendor-proto/tmp && \
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p vendor-proto/validate
	mv vendor-proto/tmp/validate vendor-proto/
	rm -rf vendor-proto/tmp

vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

vendor-proto: vendor-rm vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api

protoc-generate:
	protoc \
	-I ${PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(BINDIR)/protoc-gen-go \
	--go_out internal/${PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(BINDIR)/protoc-gen-go-grpc \
	--go-grpc_out internal/${PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(BINDIR)/protoc-gen-validate \
    --validate_out="lang=go,paths=source_relative:internal/"${PROTO_PATH} \
    --plugin=protoc-gen-grpc-gateway=$(BINDIR)/protoc-gen-grpc-gateway \
    --grpc-gateway_out internal/${PROTO_PATH} \
    --grpc-gateway_opt logtostderr=true  \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    --plugin=protoc-gen-openapiv2=$(BINDIR)/protoc-gen-openapiv2 \
    --openapiv2_out internal/${PROTO_PATH} \
    --openapiv2_opt logtostderr=true \
	api/loms.proto && \
	go mod tidy

install-goose:
	GOBIN=$(BINDIR) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

db-create-migration:
	mkdir -p  migrations
	$(BINDIR)/goose -dir $(MIGRATIONS_FOLDER) create -s $(n) add_loms_table sql
	$(BINDIR)/goose -dir $(MIGRATIONS_FOLDER) create -s $(n) init_loms_table sql

db-migrate:
	$(BINDIR)/goose -dir $(MIGRATIONS_FOLDER) postgres "$(LOCAL_DB_DSN)" up

db-migrate-down:
	$(BINDIR)/goose -dir $(MIGRATIONS_FOLDER) postgres "$(LOCAL_DB_DSN)" down



