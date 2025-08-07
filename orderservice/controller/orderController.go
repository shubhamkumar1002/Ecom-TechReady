package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"log"
	"orderService/model"
	"orderService/service"
)

type OrderController struct {
	Service service.OrderService
}

func (oc *OrderController) CreateOrder(ctx iris.Context) {
	var order model.OrderCreateDTO
	if err := ctx.ReadJSON(&order); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	newOrder, err := oc.Service.Create(&order)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create order"})
		log.Printf("err details :   ", err)
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(newOrder)
}

func (oc *OrderController) GetOrderByID(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid order ID"})
		return
	}

	order, err := oc.Service.GetOrderById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": fmt.Sprintf("Order not found with OrderID: %s", id)})
		return
	}

	ctx.JSON(order)
}

func (oc *OrderController) GetOrders(ctx iris.Context) {
	result, err := oc.Service.GetOrders()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Internal server error"})
	}

	if result == nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "No orders found"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Orders retrieved successfully", "orders": result})
}

func (oc *OrderController) UpdateOrderStatus(ctx iris.Context) {
	var orderupdate model.OrderUpdateDTO
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid order ID"})
		return
	}

	if err := ctx.ReadJSON(&orderupdate); err != nil || (orderupdate.OrderStatus == "" && orderupdate.PaymentStatus == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid status"})
		return
	}

	if _, err := oc.Service.UpdateOrderStatus(id, orderupdate.OrderStatus, orderupdate.PaymentStatus); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to update status"})
		return
	}

	ctx.JSON(iris.Map{"message": "Order status updated"})
}
