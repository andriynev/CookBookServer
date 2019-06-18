package dictionary

type Unit struct {
	Dictionary
}

func (Unit) TableName() string {
	return "units"
}
