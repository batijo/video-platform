package models

import "github.com/jinzhu/gorm"

//User struct declaration
type User struct {
	gorm.Model
	Name     string `json:"Name"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Gender   string `json:"Gender"`
	Password string `json:"Password"`

	Video  []Video   `gorm:"ForeignKey:UserID"`
	Stream []Vstream `gorm:"ForeignKey:UserID"`
}
