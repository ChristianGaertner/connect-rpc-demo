package order

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	pb "github.com/ChristianGaertner/connect-rpc-demo/proto/order/v1"
	"github.com/ChristianGaertner/connect-rpc-demo/proto/order/v1/orderv1connect"
	productv1 "github.com/ChristianGaertner/connect-rpc-demo/proto/product/v1"
	"github.com/ChristianGaertner/connect-rpc-demo/proto/product/v1/productv1connect"
)

type Service struct {
	ps productv1connect.ProductServiceClient
}

func Register(ctx context.Context, mux *http.ServeMux, ps productv1connect.ProductServiceClient) error {
	srvc := Service{
		ps: ps,
	}

	validationInterceptor, err := validate.NewInterceptor()
	if err != nil {
		return fmt.Errorf("create valdidate interceptor: %w", err)
	}

	mux.Handle(orderv1connect.NewOrderServiceHandler(srvc, connect.WithInterceptors(
		validationInterceptor,
	)))
	return nil
}

func (s Service) CreateOrder(ctx context.Context, req *connect.Request[pb.CreateOrderRequest]) (*connect.Response[pb.CreateOrderResponse], error) {
	fmt.Println(req.Msg.GetId())
	for _, l := range req.Msg.GetItems() {
		productReq := connect.NewRequest(&productv1.GetProductRequest{Id: l.GetProductId()})
		resp, err := s.ps.GetProduct(ctx, productReq)
		if err != nil {
			return nil, fmt.Errorf("get product: %w", err)
		}
		fmt.Println("Product:", resp.Msg.GetName())
	}
	return connect.NewResponse(&pb.CreateOrderResponse{}), nil
}
