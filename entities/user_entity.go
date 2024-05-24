package entities

type UserEntity struct {
	IdUser     uint  `gorm:"primarykey"`
	IdTelegram int64 `gorm:"unique;not null"`
	Admin      bool  `gorm:"default:false"`
}

func (UserEntity) TableName() string {
	return "users"
}
