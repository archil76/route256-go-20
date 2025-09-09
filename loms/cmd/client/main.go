package main

import (
	"context"
	"fmt"
	desc "route256/loms/internal/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		":8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := desc.NewLomsClient(conn)
	ctx := context.Background()

	resp, err := client.StocksInfo(ctx, &desc.StocksInfoRequest{Sku: 139275865})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
