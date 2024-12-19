package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primaryKey; autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

// migration is required to start db in case of postgress
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
