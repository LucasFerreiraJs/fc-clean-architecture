package usecase

import (
	"github.com/devfullcycle/fc-clean-architecture/internal/entity"
)

type GetOrdersOutput struct {
	ID         string  `json:"id"`
	FinalPrice float64 `json:"final_price"`
}

type GetOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *GetOrdersUseCase {
	return &GetOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *GetOrdersUseCase) Execute() ([]GetOrdersOutput, error) {

	output, err := l.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// var listOrder []OrderOutputDTO
	var listOrder []GetOrdersOutput
	for _, order := range output {
		listOrder = append(listOrder, GetOrdersOutput{
			ID:         order.ID,
			FinalPrice: order.FinalPrice,
		})
	}

	return listOrder, nil

}
