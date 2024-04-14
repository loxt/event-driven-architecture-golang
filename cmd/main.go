package main

import (
	"context"
	"fmt"
	"github.com/loxt/event-driven-architecture-golang/internal/application/controller"
	"github.com/loxt/event-driven-architecture-golang/internal/application/usecase"
	"github.com/loxt/event-driven-architecture-golang/internal/domain/event"
	"github.com/loxt/event-driven-architecture-golang/internal/infra/queue"
	"log"
	"net/http"
	"reflect"
)

func main() {
	ctx := context.Background()

	// initialize memoryQueueAdapter
	memoryQueueAdapter := queue.NewMemoryQueueAdapter()

	// use cases
	createOrderUseCase := usecase.NewCreateOrderUseCase(memoryQueueAdapter)
	processPaymentUseCase := usecase.NewProcessOrderPaymentUseCase(memoryQueueAdapter)
	stockMovementUseCase := usecase.NewStockMovementUseCase()
	sendOrderEmailUseCase := usecase.NewSendOrderEmailUseCase()

	// controllers
	orderController := controller.NewOrderController(createOrderUseCase, processPaymentUseCase, stockMovementUseCase, sendOrderEmailUseCase)

	// register routes
	http.HandleFunc("POST /create-order", orderController.CreateOrder)

	// mapping listeners
	var list map[reflect.Type][]func(w http.ResponseWriter, r *http.Request) = map[reflect.Type][]func(w http.ResponseWriter, r *http.Request){
		reflect.TypeOf(event.OrderCreatedEvent{}): {
			orderController.ProcessOrderPayment,
			orderController.StockMovement,
			orderController.SendOrderEmail,
		},
	}

	// register listeners
	for eventType, handlers := range list {
		for _, handler := range handlers {
			memoryQueueAdapter.ListenerRegister(eventType, handler)
		}
	}

	// connect memoryQueueAdapter
	err := memoryQueueAdapter.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connect memoryQueueAdapter %s", err)
	}
	defer memoryQueueAdapter.Disconnect(ctx)

	// start consuming queues
	OrderCreatedEvent := reflect.TypeOf(event.OrderCreatedEvent{}).Name()

	go func(ctx context.Context, queueName string) {
		err = memoryQueueAdapter.StartConsuming(ctx, queueName)
		if err != nil {
			log.Fatalf("Error running consumer %s: %s", queueName, err)
		}
	}(ctx, OrderCreatedEvent)

	// start server
	fmt.Println("Server is running on port 8080")
	_ = http.ListenAndServe(":8080", nil)
}
