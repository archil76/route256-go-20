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
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

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

vendor-proto: vendor-rm vendor-proto/google/protobuf vendor-proto/validate

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
	api/loms.proto && \
	go mod tidy


#protoc-generate: bin-deps
#	protoc \
#	-I ./api \
#	--go_out=./internal/api \
#	--go_opt=paths=source_relative \
#	--go-grpc_out=./internal/api \
#	--go-grpc_opt=paths=source_relative \
#    --validate_out="lang=go,paths=source_relative:api" \
#	api/loms.proto && \
#	go mod tidy