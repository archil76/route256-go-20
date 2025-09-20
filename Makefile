include make/lint.mk
include make/build.mk

lint: cart-lint loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

run:
	export CONFIG_FILE=./loms/configs/values_local.yaml &&\
	go run ./loms/cmd/server/main.go

run-all: bindir
	echo "build cart"
	export CONFIG_FILE=./cart/configs/values_local.yaml &&\
    go run ./cart/cmd/server/main.go
	echo "build loms"
	export CONFIG_FILE=./loms/configs/values_local.yaml &&\
    go run ./loms/cmd/server/main.go
