package main

import (
	"fmt"
	"github.com/loxt/event-driven-architecture-golang/internal/application/controller"
	"github.com/loxt/event-driven-architecture-golang/internal/application/usecase"
	"github.com/loxt/event-driven-architecture-golang/internal/infra/queue"
	"net/http"
)

func main() {
	memoryQueueAdapter := queue.NewMemoryQueueAdapter()
	createOrderUseCase := usecase.NewCreateOrderUseCase(memoryQueueAdapter)
	orderController := controller.NewOrderController(createOrderUseCase)

	http.HandleFunc("POST /create-order", orderController.CreateOrder)

	fmt.Println("Server is running on port 8080")
	_ = http.ListenAndServe(":8080", nil)
}
