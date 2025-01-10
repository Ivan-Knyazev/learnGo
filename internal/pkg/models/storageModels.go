package models

// Main Model for Data
type Data struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Key       string  `gorm:"not null" json:"key"`
	ValueType string  `gorm:"not null" json:"valueType"`
	ExpiresAt int64   `gorm:"not null" json:"expiresAt"`
	Scalar    Scalar  `gorm:"foreignKey:DataID;constraint:OnDelete:CASCADE;" json:"scalar"`
	Slice     []Slice `gorm:"foreignKey:DataID;constraint:OnDelete:CASCADE;" json:"slice"`
	Dict      []Dict  `gorm:"foreignKey:DataID;constraint:OnDelete:CASCADE;" json:"dict"`
}

// Model for ValueType Scalar
type Scalar struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	DataID            uint   `gorm:"not null" json:"dataID"`
	ScalarValueType   string `gorm:"not null" json:"scalarValueType"`
	ScalarValueInt    int    `json:"scalarValueInt"`
	ScalarValueString string `json:"scalarValueString"`
}

// Model for ValueType Slice
type Slice struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	DataID uint `gorm:"not null" json:"dataID"`
	Value  int  `gorm:"not null" json:"value"`
}

// Model for ValueType Dict
type Dict struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	DataID            uint   `gorm:"not null" json:"dataID"`
	Field             string `gorm:"not null" json:"field"`
	ScalarValueType   string `gorm:"not null" json:"scalarValueType"`
	ScalarValueInt    int    `json:"scalarValueInt"`
	ScalarValueString string `json:"scalarValueString"`
}
