package dictionary

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type DictEntity interface {
	TableName() string
	SetValue(value string)
}

type DictionaryRepository struct {
	db *gorm.DB
}

func GetDictionaryRepository(db *gorm.DB) *DictionaryRepository {
	return &DictionaryRepository{db: db}
}

func (r *DictionaryRepository) GetByValue(value string, entity DictEntity) (err error) {
	if len(value) == 0 {
		err = fmt.Errorf("param value cannot be empty")
		return
	}
	err = r.db.Table(entity.TableName()).Where(Dictionary{Value: value}).First(entity).Error
	return
}

func (r *DictionaryRepository) GetById(id uint, entity DictEntity) (err error) {
	if id == 0 {
		err = fmt.Errorf("param id cannot be empty")
		return
	}
	err = r.db.Table(entity.TableName()).Where(Dictionary{Id: id}).First(entity).Error
	return
}

func (r *DictionaryRepository) GetByIds(ids []uint, entity DictEntity) (items []Dictionary, err error) {
	if len(ids) == 0 {
		return
	}
	err = r.db.Table(entity.TableName()).Where("id in (?)", ids).Find(&items).Error
	return
}

func (r *DictionaryRepository) GetAll(entity DictEntity) (items []Dictionary, err error) {
	err = r.db.Table(entity.TableName()).Find(&items).Error
	return
}

func (r *DictionaryRepository) Create(value string, entity DictEntity) (err error) {
	if len(value) == 0 {
		err = fmt.Errorf("param value cannot be empty")
		return
	}
	entity.SetValue(value)
	err = r.db.Table(entity.TableName()).Create(entity).Error
	return
}