package models

type Order struct {
	OrderId      uint   `json:"orderId" gorm:"primaryKey"`
	CustomerName string `json:"customerName" gorm:"not null;type:varchar(100)"`
	OrderedAt    string `json:"orderedAt" gorm:"not null;type:timestamp"`
	Items        []Item `json:"items"`
}