package service

import (
	"context"

	"github.com/devfullcycle/fc-clean-architecture/internal/infra/grpc/pb"
	"github.com/devfullcycle/fc-clean-architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	GetOrdersUseCase   usecase.GetOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, getOrdersUseCase usecase.GetOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		GetOrdersUseCase:   getOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) GetOrders(ctx context.Context, in *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {

	orders, err := s.GetOrdersUseCase.Execute()

	if err != nil {
		return nil, err
	}

	var response []*pb.CreateOrderResponse

	for _, order := range orders {
		orderItem := &pb.CreateOrderResponse{
			Id:         order.ID,
			FinalPrice: float32(order.FinalPrice),
		}
		// Price:      float32(order.Price),
		// Tax:        float32(order.Price),
		response = append(response, orderItem)
	}

	return &pb.GetOrdersResponse{
		OrderList: response,
	}, nil
}
