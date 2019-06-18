package receipt

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
}

func GetReceiptRepository(db *gorm.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) GetAll() (receipts []Receipt, err error) {
	err = r.db.Preload("Media").Find(&receipts).Error
	return
}

func (r *ReceiptRepository) GetAllByCategory(category string) (receipts []Receipt, err error) {
	if category == "" {
		receipts, err = r.GetAll()
		return
	}
	err = r.db.Preload("Media").Where(&Receipt{Category:category}).Find(&receipts).Error
	return
}

func (r *ReceiptRepository) GetById(id uint) (receipt Receipt, err error) {
	if id == 0 {
		err = fmt.Errorf("receipt id cannot be empty")
		return
	}
	err = r.db.Where(&Receipt{Id: id}).
		Preload("Media").
		First(&receipt).Error
	return
}

func (r *ReceiptRepository) CreateIngredient(receiptIngredient *ReceiptIngredient) (err error) {
	err = r.db.Create(receiptIngredient).Error
	return
}

func (r *ReceiptRepository) UpdateIngredient(receiptIngredient *ReceiptIngredient) (err error) {
	err = r.db.Model(&ReceiptIngredient{}).Where(&ReceiptIngredient{Id: receiptIngredient.Id}).Update(receiptIngredient).Error
	return
}

func (r *ReceiptRepository) GetIngredientById(id uint) (ingredient ReceiptIngredient, err error) {
	if id == 0 {
		err = fmt.Errorf("receipt id cannot be empty")
		return
	}

	err = r.db.Model(ReceiptIngredient{}).Preload("Ingredient").Where(&ReceiptIngredient{Id: id}).First(&ingredient).Error
	return
}

func (r *ReceiptRepository) DeleteIngredientById(id uint) (err error) {
	if id == 0 {
		err = fmt.Errorf("id cannot be empty")
		return
	}

	err = r.db.Model(ReceiptIngredient{}).Where(&ReceiptIngredient{Id: id}).Delete(&ReceiptIngredient{}).Error
	return
}

func (r *ReceiptRepository) CreateDirection(direction *ReceiptDirection) (err error) {
	err = r.db.Create(direction).Error
	return
}

func (r *ReceiptRepository) UpdateDirection(direction *ReceiptDirection) (err error) {
	err = r.db.Model(&ReceiptDirection{}).Where(&ReceiptDirection{Id: direction.Id}).Update(direction).Error
	return
}

func (r *ReceiptRepository) GetDirectionById(id uint) (direction ReceiptDirection, err error) {
	if id == 0 {
		err = fmt.Errorf("receipt id cannot be empty")
		return
	}

	err = r.db.Model(ReceiptDirection{}).Where(&ReceiptDirection{Id: id}).First(&direction).Error
	return
}

func (r *ReceiptRepository) DeleteDirectionById(id uint) (err error) {
	if id == 0 {
		err = fmt.Errorf("id cannot be empty")
		return
	}

	err = r.db.Model(ReceiptDirection{}).Where(&ReceiptDirection{Id: id}).Delete(&ReceiptDirection{}).Error
	return
}

func (r *ReceiptRepository) GetIngredientsById(id uint) (ingredients []ReceiptIngredient, err error) {
	if id == 0 {
		err = fmt.Errorf("receipt id cannot be empty")
		return
	}

	err = r.db.Model(ReceiptIngredient{}).Preload("Ingredient").Where(&ReceiptIngredient{ReceiptId: id}).Find(&ingredients).Error
	return
}



func (r *ReceiptRepository) GetDirectionsById(id uint) (directions []ReceiptDirection, err error) {
	if id == 0 {
		err = fmt.Errorf("receipt id cannot be empty")
		return
	}

	err = r.db.Model(ReceiptIngredient{}).Where(&ReceiptDirection{ReceiptId: id}).Find(&directions).Error
	return
}

func (r *ReceiptRepository) Create(receipt *Receipt) (err error) {
	if receipt == nil {
		err = fmt.Errorf("receipt cannot be empty")
		return
	}
	if receipt.Id != 0 {
		err = fmt.Errorf("receipt id should be empty")
		return
	}
	err = r.db.Create(receipt).Error
	return
}

func (r *ReceiptRepository) Update(receipt *Receipt) (err error) {
	if receipt == nil {
		err = fmt.Errorf("receipt cannot be empty")
		return
	}
	if receipt.Id == 0 {
		err = fmt.Errorf("receipt id cannot be empty")
		return
	}

	err = r.db.Model(Receipt{}).Where(&Receipt{Id: receipt.Id}).Update(receipt).Error
	return
}

func (r *ReceiptRepository) Delete(id uint) (err error) {
	if id == 0 {
		err = fmt.Errorf("media id cannot be empty")
		return
	}

	err = r.db.Model(Receipt{}).Where(&Receipt{Id: id}).Delete(Receipt{}).Error
	return
}
