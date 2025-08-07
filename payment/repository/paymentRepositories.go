package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"paymentService/model"
	"time"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (paym *PaymentRepository) CheckOrder(orderid uuid.UUID) (bool, error) {
	var payment model.Payment
	err := paym.DB.Model(&payment).Where("order_id = ?", orderid).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (paym *PaymentRepository) CreatePayment(paymentCreateDTO *model.PaymentCreateDTO) (*model.Payment, error) {
	paymentId := uuid.New()

	newPayment := &model.Payment{
		ID:            paymentId,
		OrderID:       paymentCreateDTO.OrderID,
		TotalAmount:   paymentCreateDTO.TotalAmount,
		PaymentStatus: paymentCreateDTO.PaymentStatus,
		CreatedAt:     time.Now(),
	}

	err := paym.DB.Create(&newPayment).Error
	if err != nil {

		return nil, err
	}

	return newPayment, nil
}

func (pay *PaymentRepository) UpdateStatus(id uuid.UUID, status string) error {
	var updatepayment model.Payment
	err := pay.DB.Model(&updatepayment).Where("order_id = ?", id).Update("payment_status", status).Update("updated_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (pay *PaymentRepository) GetPaymentByOrderId(orderid uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	if err := pay.DB.Where("order_id", orderid).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (pay *PaymentRepository) GetPayments() ([]model.Payment, error) {
	var payments []model.Payment
	err := pay.DB.Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
