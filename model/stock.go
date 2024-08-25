package model

type Stock struct {
	ID              int `gorm:"primaryKey"`
	Milk            int `gorm:"not null"`
	MilkLactoseFree int `gorm:"default:0"`
	Arabica         int `gorm:"default:0"`
	Robusta         int `gorm:"default:0"`
	Espresso        int `gorm:"default:0"`
	Mugs            int `gorm:"default:0"`
	VanillaSyrup    int `gorm:"default:0"`
	CaramelSyrup    int `gorm:"default:0"`
	NutSyrup        int `gorm:"default:0"`
	Cream           int `gorm:"default:0"`
}
