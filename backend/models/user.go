package models

import "github.com/jinzhu/gorm"

//User struct declaration
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`
	Role     string `json:"role"`

	Video  []Video   `gorm:"ForeignKey:UserID"`
	Stream []Vstream `gorm:"ForeignKey:UserID"`
}
