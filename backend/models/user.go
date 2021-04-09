package models

import "github.com/jinzhu/gorm"

//User struct declaration
type User struct {
	gorm.Model
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
	Public   bool   `json:"public"`

	Video  []Video   `gorm:"ForeignKey:UserID"`
	Stream []Vstream `gorm:"ForeignKey:UserID"`
}
