package models

type Wallet struct {
	Id      int64 `gorm:"primaryKey" json:"id"`
	UserId  int   `gorm:"type:Int" json:"user_id"`
	Balance int   `gorm:"type:Int;default:0" json:"balance"`

	IsDeleted int `gorm:"default:0" json:"is_deleted"`
	CreatedAt int `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int `gorm:"autoUpdateTime:milli" json:"updated_at"`
	CreatedBy int `json:"created_by"`
	UpdatedBy int `json:"updated_by"`
}
