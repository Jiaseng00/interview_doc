package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	Id        int    `gorm:"id" json:"id"`
	Name      string `gorm:"name" json:"name"`
	Email     string `gorm:"email" json:"email"`
	CreatedAt string `gorm:"created_at" json:"created_at"`
}

func main() {
	// 模拟链接MySQL数据库
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}

	// 链接完毕
	var users []User
	sevenDaysAgoTime := time.Now().AddDate(0, 0, -7).Unix()
	err = db.Table("users").Where("created_at >= ?", sevenDaysAgoTime).Find(&users).Error
	if err != nil {
		log.Println(err)
		return
	}
}
