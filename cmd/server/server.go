package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ChristianGaertner/connect-rpc-demo/internal/order"
	"github.com/ChristianGaertner/connect-rpc-demo/internal/product"
	"github.com/ChristianGaertner/connect-rpc-demo/proto/product/v1/productv1connect"
)

var (
	serviceFlag = flag.String("service", "", "Service to run (order or product)")
)

var dns = map[string]string{
	"order":   "localhost:7444",
	"product": "localhost:7445",
}

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()

	var err error
	switch *serviceFlag {
	case "order":
		ps := productv1connect.NewProductServiceClient(http.DefaultClient, "http://"+dns["product"])
		err = order.Register(ctx, mux, ps)
	case "product":
		err = product.Register(ctx, mux)
	default:
		err = fmt.Errorf("unknown service: %q", *serviceFlag)
	}
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}
	addr := dns[*serviceFlag]
	fmt.Println(*serviceFlag, "service started on", addr)
	return http.ListenAndServe(addr, mux)
}
