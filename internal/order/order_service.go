package order

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	pb "github.com/ChristianGaertner/connect-rpc-demo/proto/order/v1"
	"github.com/ChristianGaertner/connect-rpc-demo/proto/order/v1/orderv1connect"
)

type Service struct{}

func Register(ctx context.Context, mux *http.ServeMux) error {
	srvc := Service{}
	mux.Handle(orderv1connect.NewOrderServiceHandler(srvc))
	return nil
}

func (Service) CreateOrder(ctx context.Context, req *connect.Request[pb.CreateOrderRequest]) (*connect.Response[pb.CreateOrderResponse], error) {
	fmt.Println(req.Msg.GetId())
	return connect.NewResponse(&pb.CreateOrderResponse{}), nil
}
