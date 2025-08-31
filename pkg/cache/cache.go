package cache

import (
	"sync"
	"wb-task-L0/pkg/models"
)

type OrderCache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]models.Order),
	}
}

func (c *OrderCache) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *OrderCache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[id]
	return order, ok
}

func (c *OrderCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.orders, id)
}

func (c *OrderCache) GetAll() []models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	orders := make([]models.Order, 0, len(c.orders))
	for _, order := range c.orders {
		orders = append(orders, order)
	}
	return orders
}

func (c *OrderCache) LoadFromDB(orders []models.Order) {
	newMap := make(map[string]models.Order, len(orders))
	for _, o := range orders {
		newMap[o.OrderUID] = o
	}
	c.mu.Lock()
	c.orders = newMap
	c.mu.Unlock()
}

func (c *OrderCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.orders)
}
