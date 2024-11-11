package models

type User struct {
	Id        int    `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"type:varchar(100);unique" json:"username"`
	Password  string `gorm:"type:varchar(100)" json:"password"`
	IsDeleted int    `gorm:"default:0" json:"is_deleted"`
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"updated_at"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}

func SeedAdminUser() {
	admin := User{Username: "gibran", Password: "gibran123"}
	DB.Create(&admin)
}
