package controller

import (
	"Petstore/models/store/entity"
	"Petstore/models/store/service"
	"Petstore/responder"
	"github.com/go-chi/chi/v5"
	"github.com/ptflp/godecoder"
	"net/http"
	"strconv"
)

type StoreController interface {
	Inventory(http.ResponseWriter, *http.Request)
	PlaceOrder(http.ResponseWriter, *http.Request)
	FindOrderById(http.ResponseWriter, *http.Request)
	DeleteById(http.ResponseWriter, *http.Request)
}

type StoreControl struct {
	service service.OrderServicer
	responder.Responder
	godecoder.Decoder
}

func NewStoreController(service service.OrderServicer, responder responder.Responder, decoder godecoder.Decoder) *StoreControl {
	return &StoreControl{service: service, Responder: responder, Decoder: decoder}
}

// Inventory Returning a map of status codes to quantities
// @Summary Returns pet inventories by status
// @Description Return pet inventories by status
// @Tags store
// @Accept json
// @Produce json
// @Success 200 {object} InventoryResponse
// @Router /store/inventory [get]
func (c *StoreControl) Inventory(w http.ResponseWriter, r *http.Request) {
	out := c.service.Inventory(r.Context())
	c.OutputJSON(w, out)
}

// PlaceOrder Placing an order for a pet
// @Summary Place an order for a pet
// @Description Returns a map of status codes to quantities
// @Tags store
// @Accept json
// @Produce json
// @Param Order body entity.Order true "order placed for purchasing the pet"
// @Success 200 {object} PlaceOrderResponse
// @Router /store/order [post]
func (c *StoreControl) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var order entity.Order
	if err := c.Decode(r.Body, &order); err != nil {
		c.ErrorBadRequest(w, err)
		return
	}
	out, err := c.service.PlaceOrder(r.Context(), order)
	if err != nil {
		c.ErrorInternal(w, err)
		return
	}

	c.OutputJSON(w, PlaceOrderResponse{
		Code:    http.StatusCreated,
		Success: true,
		Data:    Data{Order: out},
	})
}

// FindOrderById Finding purchased order by ID
// @Summary Find purchased order by ID
// @Description Returning order by id
// @Tags store
// @Accept  json
// @Produce  json
// @Param   id path int true "ID of order to return"
// @Success 200 {object} PlaceOrderResponse
// @Router /store/order/{id} [get]
func (c *StoreControl) FindOrderById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, _ := strconv.Atoi(id)

	out, err := c.service.FindOrderById(r.Context(), uint64(intID))
	if err != nil {
		c.ErrorInternal(w, err)
		return
	}

	c.OutputJSON(w, PlaceOrderResponse{
		Code:    http.StatusOK,
		Success: true,
		Data:    Data{Order: out},
	})

}

// DeleteById Deleting purchased order by ID
// @Summary Delete purchased order by ID
// @Description Deleting purchased order by ID
// @Tags store
// @Accept  json
// @Produce  json
// @Param   id path int true "ID of order to delete"
// @Success 200 {object} DeleteOrderResponse
// @Router /store/order/{id} [delete]
func (c *StoreControl) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, _ := strconv.Atoi(id)
	err := c.service.DeleteById(r.Context(), uint64(intID))
	if err != nil {
		c.OutputJSON(w, DeleteOrderResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Order Not Found",
		})
		return
	}

	c.OutputJSON(w, DeleteOrderResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Order Deleted Successfully",
	})

}
