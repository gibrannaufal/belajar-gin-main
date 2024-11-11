package models

type Product struct {
	Id           int64  `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(100)" json:"name"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
	Price        int    `json:"price"`
	Availability int    `json:"availability"`

	IsDeleted int `gorm:"default:0" json:"is_deleted"`
	CreatedAt int `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int `gorm:"autoUpdateTime:milli" json:"updated_at"`
	CreatedBy int `json:"created_by"`
	UpdatedBy int `json:"updated_by"`
}
