package models

import "gorm.io/gorm"

type Users struct {
	ID       uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	BookID   uint   `json:"bookId"`
	Book     Books  `gorm:"foreignKey:BookID" json:"book"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
