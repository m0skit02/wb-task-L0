package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wb-task-L0/pkg/models"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *models.Order) (string, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return order.OrderUID, nil
}

func (r *OrderRepo) CreateOrderWithAssociations(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Проверяем, есть ли уже заказ
		var existing models.Order
		if err := tx.Where("order_uid = ?", order.OrderUID).First(&existing).Error; err == nil {
			// заказ уже есть, пропускаем вставку
			return nil
		}

		// Вставляем заказ
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(order).Error; err != nil {
			return err
		}

		// Вставляем Delivery
		order.Delivery.DeliveryID = order.Delivery.DeliveryID + "_" + order.OrderUID // уникальный ключ на всякий случай
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&order.Delivery).Error; err != nil {
			return err
		}

		// Вставляем Payment
		order.Payment.PaymentID = order.Payment.PaymentID + "_" + order.OrderUID
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&order.Payment).Error; err != nil {
			return err
		}

		// Вставляем Items
		for i := range order.Items {
			order.Items[i].ItemID = order.Items[i].ItemID + "_" + order.OrderUID
		}
		if len(order.Items) > 0 {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&order.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *OrderRepo) GetAll() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepo) GetByID(orderUID string) (models.Order, error) {
	var order models.Order

	// Ищем заказ и сразу подгружаем ассоциации
	if err := r.db.
		Preload("Delivery").
		Preload("Payment").
		Preload("Items").
		First(&order, "order_uid = ?", orderUID).Error; err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *OrderRepo) Delete(orderUID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Удаляем связанные items
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Item{}).Error; err != nil {
			return err
		}

		// Удаляем payment
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Payment{}).Error; err != nil {
			return err
		}

		// Удаляем delivery
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Delivery{}).Error; err != nil {
			return err
		}

		// Удаляем сам заказ
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Order{}).Error; err != nil {
			return err
		}

		return nil
	})
}
