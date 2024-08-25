package model

type Order struct {
	ID        uint       `gorm:"primaryKey"`
	Beverages []Beverage `gorm:"many2many:order_beverages;"`
	Completed bool       `gorm:"default:false"`
	Price     float32
}
