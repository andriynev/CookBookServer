package media

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	MediaFolderRoot = "/data/media/"

	Unknown Type = iota
	Image
	Video
)

type Type uint

func (t Type) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *Type) Scan(src interface{}) error {

	switch orig := src.(type) {
	case int32:
		*t = Type(orig)
	case int64:
		*t = Type(orig)
	default:
		return fmt.Errorf("incompatible type for Type %s", orig)
	}

	return nil
}

type Media struct {
	Id        uint      `json:"-" gorm:"primary_key"`
	// media link (read only) (required)
	Link      string    `json:"link" example:"ac146571eb55d9cbde1a886dbc1f85d2/event.png"`
	// format (read only) (required)
	Format    string    `json:"format" example:"image/png"`
	CreatedAt time.Time `json:"-"`
}

func (Media) TableName() string {
	return "media"
}
