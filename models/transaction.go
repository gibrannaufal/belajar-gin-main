package models

type Transaction struct {
	Id              int    `gorm:"primaryKey" json:"id"`
	WalletId        int    `gorm:"type:Int" json:"wallet_id"`
	TransactionType string `gorm:"type:enum('withdraw','deposit');default:NULL" json:"transaction_type"`
	Amount          int    `gorm:"type:Int" json:"amout"`

	IsDeleted int `gorm:"default:0" json:"is_deleted"`
	CreatedAt int `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int `gorm:"autoUpdateTime:milli" json:"updated_at"`
	CreatedBy int `json:"created_by"`
	UpdatedBy int `json:"updated_by"`
}
