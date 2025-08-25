FROM golang:1.23.9-alpine as builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

FROM scratch
COPY --from=builder server /bin/server
COPY configs/values_local.yaml /bin/config/values_local.yaml

ENV CONFIG_FILE=/bin/config/values_local.yaml

ENTRYPOINT ["/bin/server"]