package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Goly struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Redirect string `json:"redirect" gorm:"not null"`
	Goly     string `json:"goly" gorm:"unique;not null"`
	Clicked  uint64 `json:"clicked"`
	Random   bool   `json:"random"`
}

func Setup() {
	dsn := "root:test@tcp(127.0.0.1:3306)/goly"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Goly{})
	if err != nil {
		fmt.Println(err)
	}
}
