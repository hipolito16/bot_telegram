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
	DB       *gorm.DB
	err      error
	host     string
	port     string
	user     string
	password string
	dbname   string
)

func Start() {
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname = os.Getenv("DB_NAME")

	stringDeConexao := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbname)

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
		adminName := os.Getenv("ADMIN_NAME")
		user := entities.UserEntity{IdTelegram: adminIdTelegram, Admin: true, Name: adminName}
		DB.Create(&user)
	}
}
