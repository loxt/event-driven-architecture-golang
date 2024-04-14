package usecase

import (
	"context"
	"fmt"
	"github.com/loxt/event-driven-architecture-golang/internal/domain/event"
)

type StockMovementUseCase struct {
}

func NewStockMovementUseCase() *StockMovementUseCase {
	return &StockMovementUseCase{}
}

func (h *StockMovementUseCase) Execute(ctx context.Context, payload *event.OrderCreatedEvent) error {
	for _, item := range payload.Items {
		fmt.Printf("Removing %d items of product %s from stock\n", item.Quantity, item.ProductName)
	}
	return nil
}
