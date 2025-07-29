package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ChristianGaertner/connect-rpc-demo/internal/order"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()

	err := order.Register(ctx, mux)
	if err != nil {
		return err
	}

	return http.ListenAndServe("localhost:7444", mux)
}
