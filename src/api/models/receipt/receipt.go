package receipt

import (
	"food/src/api/models/ingredient"
	"food/src/api/models/media"
	"time"
)

type Receipt struct {
	Id        uint      `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Description string `json:"description"`
	Category string `json:"category"`
	CookingTime int `json:"cooking_time"`
	UserId uint `json:"user_id"`
	MediaId *uint `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Media *media.Media `gorm:"foreignkey:MediaId" json:"media,omitempty"`
}

func (Receipt) TableName() string {
	return "receipts"
}

type ReceiptIngredient struct {
	Id        uint      `json:"id" gorm:"primary_key"`
	Quantity string `json:"quantity"`
	ReceiptId uint `json:"receipt_id"`
	IngredientId uint `json:"ingredient_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Ingredient *ingredient.Ingredient `json:"ingredient,omitempty" gorm:"foreignkey:IngredientId"`
}

func (ReceiptIngredient) TableName() string {
	return "receipt_ingredients"
}

type ReceiptDirection struct {
	Id        uint      `json:"id" gorm:"primary_key"`
	ReceiptId uint `json:"receipt_id"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

func (ReceiptDirection) TableName() string {
	return "receipt_directions"
}
