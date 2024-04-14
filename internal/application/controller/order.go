package controller

import (
	"encoding/json"
	"github.com/loxt/event-driven-architecture-golang/internal/application/dto"
	"github.com/loxt/event-driven-architecture-golang/internal/application/usecase"
	"net/http"
)

type OrderController struct {
	createOrderUserCase *usecase.CreateOrderUseCase
}

func NewOrderController(createOrderUserCase *usecase.CreateOrderUseCase) *OrderController {
	return &OrderController{
		createOrderUserCase: createOrderUserCase,
	}
}

func (u *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var requestData dto.CreateOrderDTO
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	err = u.createOrderUserCase.Execute(r.Context(), requestData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
