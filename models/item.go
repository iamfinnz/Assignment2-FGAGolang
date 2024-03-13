package models

type Item struct {
	ItemId      uint   `json:"lineItemId" gorm:"primaryKey"`
	ItemCode    string `json:"itemCode" gorm:"not null;type:varchar(20)"`
	Description string `json:"description" gorm:"not null"`
	Quantity    uint   `json:"quantity" gorm:"not null"`
	OrderId     uint   `json:"orderId"`
}