package dictionary

import "time"

type Dictionary struct {
	Id          uint       `json:"id" gorm:"primary_key"`
	Value string `json:"value"`
	CreatedAt   time.Time  `json:"-"`
}

func (d *Dictionary) SetValue(value string) {
	d.Value = value
}

type DictionaryItems []Dictionary

func (i DictionaryItems) ConvertToMap() map[uint]string {
	items := make(map[uint]string, len(i))

	for _, item := range i {
		items[item.Id] = item.Value
	}
	return items
}