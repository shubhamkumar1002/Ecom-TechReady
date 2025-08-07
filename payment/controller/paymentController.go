package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"paymentService/service"
)

type PaymentController struct {
	Service service.PaymentService
}

func (oc *PaymentController) GetPaymentByOrderID(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid order ID"})
		return
	}

	order, err := oc.Service.GetPaymentByOrderId(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": fmt.Sprintf("payment not found with OrderID: %s", id)})
		return
	}

	ctx.JSON(order)
}

func (oc *PaymentController) GetPayments(ctx iris.Context) {
	result, err := oc.Service.GetPayments()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Internal server error"})
	}

	if result == nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "No payments found"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "payments retrieved successfully", "payments": result})
}
