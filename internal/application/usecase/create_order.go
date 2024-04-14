package usecase

import (
	"context"
	"github.com/loxt/event-driven-architecture-golang/internal/application/dto"
	"github.com/loxt/event-driven-architecture-golang/internal/domain/entity"
	"github.com/loxt/event-driven-architecture-golang/internal/domain/event"
	"github.com/loxt/event-driven-architecture-golang/internal/domain/queue"
)

type CreateOrderUseCase struct {
	publisher queue.Publisher
}

func NewCreateOrderUseCase(publisher queue.Publisher) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		publisher: publisher,
	}
}

func (u *CreateOrderUseCase) Execute(ctx context.Context, input dto.CreateOrderDTO) error {
	order, err := entity.NewOrderEntity()
	if err != nil {
		return err
	}

	for _, item := range input.Items {
		fakeProductName := "Product " + item.ProductId
		fakeProductPrice := 10.50

		i := entity.NewOrderItemEntity(fakeProductName, fakeProductPrice, item.Qtd)

		order.AddItem(i)
	}

	var eventItems []event.OrderItem
	for _, item := range order.GetItems() {
		eventItems = append(eventItems, event.OrderItem{
			ProductName: item.GetProductName(),
			TotalPrice:  item.GetTotalPrice(),
			Quantity:    item.GetQuantity(),
		})
	}

	err = u.publisher.Publish(ctx, event.OrderCreatedEvent{
		Id:         order.GetID(),
		TotalPrice: order.GetTotalPrice(),
		Status:     order.GetStatus(),
		Items:      eventItems,
	})
	if err != nil {
		return err
	}
	return nil
}
