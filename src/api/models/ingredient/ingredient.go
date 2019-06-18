package ingredient

import "time"

type Ingredient struct {
	Id        uint      `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"-"`
}

func (Ingredient) TableName() string {
	return "ingredients"
}
