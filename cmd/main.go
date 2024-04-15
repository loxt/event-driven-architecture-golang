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

	rabbitmqQueueAdapter := queue.NewRabbitMQAdapter("amqp://admin:admin@localhost:5672/")

	// use cases
	createOrderUseCase := usecase.NewCreateOrderUseCase(rabbitmqQueueAdapter)
	processPaymentUseCase := usecase.NewProcessOrderPaymentUseCase(rabbitmqQueueAdapter)
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
			rabbitmqQueueAdapter.ListenerRegister(eventType, handler)
		}
	}

	// connect queue
	err := rabbitmqQueueAdapter.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connect queue %s", err)
	}
	defer rabbitmqQueueAdapter.Disconnect(ctx)

	// start consuming queues
	OrderCreatedEvent := reflect.TypeOf(event.OrderCreatedEvent{}).Name()

	go func(ctx context.Context, queueName string) {
		err = rabbitmqQueueAdapter.StartConsuming(ctx, queueName)
		if err != nil {
			log.Fatalf("Error running consumer %s: %s", queueName, err)
		}
	}(ctx, OrderCreatedEvent)

	// start server
	fmt.Println("Server is running on port 8080")
	_ = http.ListenAndServe(":8080", nil)
}
