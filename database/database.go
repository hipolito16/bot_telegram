package database

import (
	"fmt"
	"github.com/hipolito16/bot_telegram/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

var (
	DB  *gorm.DB
	err error
)

func Start() {
	stringDeConexao := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	if DB, err = gorm.Open(postgres.Open(stringDeConexao), &gorm.Config{}); err != nil {
		panic(err)
	}

	if err = DB.AutoMigrate(&entities.UserEntity{}); err != nil {
		panic(err)
	}

	var users entities.UserEntity
	adminIdTelegram, _ := strconv.ParseInt(os.Getenv("ADMIN_ID_TELEGRAM"), 10, 64)
	tx := DB.First(&users, "id_telegram = ?", adminIdTelegram)
	if tx.RowsAffected == 0 {
		user := entities.UserEntity{IdTelegram: adminIdTelegram, Admin: true}
		DB.Create(&user)
	}
}
