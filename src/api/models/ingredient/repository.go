package ingredient

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type IngredientRepository struct {
	db *gorm.DB
}

func GetMediaRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

func (r *IngredientRepository) GetById(id uint) (ingredient Ingredient, err error) {
	if id == 0 {
		err = fmt.Errorf("ingredient id cannot be empty")
		return
	}
	err = r.db.Where(&Ingredient{Id: id}).First(&ingredient).Error
	return
}

func (r *IngredientRepository) GetAll() (ingredients []Ingredient, err error) {
	err = r.db.Find(&ingredients).Error
	return
}

func (r *IngredientRepository) Create(ingredient *Ingredient) (err error) {
	if ingredient == nil {
		err = fmt.Errorf("ingredient cannot be empty")
		return
	}
	if ingredient.Id != 0 {
		err = fmt.Errorf("ingredient id should be empty")
		return
	}
	err = r.db.Create(ingredient).Error
	return
}

func (r *IngredientRepository) Update(ingredient *Ingredient) (err error) {
	if ingredient == nil {
		err = fmt.Errorf("ingredient cannot be empty")
		return
	}
	if ingredient.Id == 0 {
		err = fmt.Errorf("ingredient id cannot be empty")
		return
	}

	err = r.db.Model(Ingredient{}).Where(&Ingredient{Id: ingredient.Id}).Update(ingredient).Error
	return
}
