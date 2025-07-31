package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ChristianGaertner/connect-rpc-demo/internal/httpmiddleware"
	"github.com/ChristianGaertner/connect-rpc-demo/internal/inproc"
	"github.com/ChristianGaertner/connect-rpc-demo/internal/order"
	"github.com/ChristianGaertner/connect-rpc-demo/internal/product"
	"github.com/ChristianGaertner/connect-rpc-demo/proto/product/v1/productv1connect"
	"golang.org/x/sync/errgroup"
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()

	ipcMux := http.NewServeMux()
	ipc := inproc.NewInMemoryServer(ipcMux)
	defer ipc.Close()

	var err error
	ps := productv1connect.NewProductServiceClient(ipc.Client(), ipc.URL())
	err = order.Register(ctx, mux, ps)
	if err != nil {
		return fmt.Errorf("register order service: %w", err)
	}
	err = product.Register(ctx, ipcMux)
	if err != nil {
		return fmt.Errorf("register product service: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		addr := "localhost:7444"
		fmt.Println(addr, "service started on", addr)
		return http.ListenAndServe(addr, httpmiddleware.CORS(mux))
	})
	return g.Wait()
}
