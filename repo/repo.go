package repo

import (
	"log"
	"uniassist/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "admin:HnVXVx8rF4G3YjS3nKuQrKVS7apg4Vzt@tcp(13.212.140.154:3306)/intern_bcc_9?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("DB failed to connect")
	}
	log.Println("database connected!")
	return db
}

var Db = connectDB() //CONNECTED TO DB

/*AFTER MIGRATION, DB WILL CLOSE ITSELF!*/
func Migration() {
	Db.AutoMigrate(&entity.User{}, &entity.Post{}, &entity.Category{}, &entity.Answer{}, &entity.LovePost{})
	log.Println("data migrated!")

	//closing db
	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatalln("DB failed to close")
	}
	sqlDB.Close()
}
