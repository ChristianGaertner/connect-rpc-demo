package product

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	pb "github.com/ChristianGaertner/connect-rpc-demo/proto/product/v1"
	"github.com/ChristianGaertner/connect-rpc-demo/proto/product/v1/productv1connect"
)

type Service struct{}

func Register(ctx context.Context, mux *http.ServeMux) error {
	srvc := Service{}

	validationInterceptor, err := validate.NewInterceptor()
	if err != nil {
		return fmt.Errorf("create valdidate interceptor: %w", err)
	}

	mux.Handle(productv1connect.NewProductServiceHandler(srvc, connect.WithInterceptors(
		validationInterceptor,
	)))
	return nil
}

func (Service) GetProduct(ctx context.Context, req *connect.Request[pb.GetProductRequest]) (*connect.Response[pb.GetProductResponse], error) {
	fmt.Println(req.Msg.GetId())
	return connect.NewResponse(&pb.GetProductResponse{
		Id:   req.Msg.GetId(),
		Name: "Product " + req.Msg.GetId(),
	}), nil
}
