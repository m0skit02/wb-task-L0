package repository

import (
	"gorm.io/gorm"
	"wb-task-L0/internal/models"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *models.Order) (string, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Сохраняем заказ вместе с зависимостями
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
