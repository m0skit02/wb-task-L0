package models

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" gorm:"column:order_uid;primaryKey"`
	TrackNumber       string    `json:"track_number" gorm:"column:track_number"`
	Entry             string    `json:"entry" gorm:"column:entry"`
	Locale            string    `json:"locale" gorm:"column:locale"`
	InternalSignature string    `json:"internal_signature" gorm:"column:internal_signature"`
	CustomerID        string    `json:"customer_id" gorm:"column:customer_id"`
	DeliveryService   string    `json:"delivery_service" gorm:"column:delivery_service"`
	ShardKey          string    `json:"shard_key" gorm:"column:shard_key"`
	SmID              int       `json:"sm_id" gorm:"column:sm_id"`
	DateCreated       time.Time `json:"date_created" gorm:"column:date_created"`
	OofShard          string    `json:"oof_shard" gorm:"column:oof_shard"`

	// Ассоциации
	Delivery Delivery `json:"delivery" gorm:"foreignKey:OrderUID;references:OrderUID"`
	Payment  Payment  `json:"payment" gorm:"foreignKey:OrderUID;references:OrderUID"`
	Items    []Item   `json:"items" gorm:"foreignKey:OrderUID;references:OrderUID"`
}

type Delivery struct {
	DeliveryID string `json:"delivery_id" gorm:"column:delivery_id;primaryKey"`
	OrderUID   string `json:"order_uid" gorm:"column:order_uid"`
	Name       string `json:"name" gorm:"column:name"`
	Phone      string `json:"phone" gorm:"column:phone"`
	Zip        string `json:"zip" gorm:"column:zip"`
	City       string `json:"city" gorm:"column:city"`
	Address    string `json:"address" gorm:"column:address"`
	Region     string `json:"region" gorm:"column:region"`
	Email      string `json:"email" gorm:"column:email"`
}

type Payment struct {
	PaymentID    string  `json:"payment_id" gorm:"column:payment_id;primaryKey"`
	OrderUID     string  `json:"order_uid" gorm:"column:order_uid"`
	Transaction  string  `json:"transaction" gorm:"column:transaction"`
	RequestID    string  `json:"request_id" gorm:"column:request_id"`
	Currency     string  `json:"currency" gorm:"column:currency"`
	Provider     string  `json:"provider" gorm:"column:provider"`
	Amount       float64 `json:"amount" gorm:"column:amount"`
	PaymentDt    int64   `json:"payment_dt" gorm:"column:payment_dt"`
	Bank         string  `json:"bank" gorm:"column:bank"`
	DeliveryCost float64 `json:"delivery_cost" gorm:"column:delivery_cost"`
	GoodsTotal   float64 `json:"goods_total" gorm:"column:goods_total"`
	CustomFee    float64 `json:"custom_fee" gorm:"column:custom_fee"`
}

type Item struct {
	ItemID      string  `json:"item_id" gorm:"column:item_id;primaryKey"`
	OrderUID    string  `json:"order_uid" gorm:"column:order_uid"`
	ChrtID      int64   `json:"chrt_id" gorm:"column:chrt_id"`
	TrackNumber string  `json:"track_number" gorm:"column:track_number"`
	Price       float64 `json:"price" gorm:"column:price"`
	Rid         string  `json:"rid" gorm:"column:rid"`
	Name        string  `json:"name" gorm:"column:name"`
	Sale        float64 `json:"sale" gorm:"column:sale"`
	Size        string  `json:"size" gorm:"column:size"`
	TotalPrice  float64 `json:"total_price" gorm:"column:total_price"`
	NmID        int64   `json:"nm_id" gorm:"column:nm_id"`
	Brand       string  `json:"brand" gorm:"column:brand"`
	Status      int     `json:"status" gorm:"column:status"`
}
