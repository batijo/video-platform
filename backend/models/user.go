package models

import "time"

//User struct declaration
type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Username string `gorm:"unique_index" json:"username"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
	Public   bool   `json:"public"`

	Video  []Video   `gorm:"ForeignKey:UserID" json:"video"`
	Stream []Vstream `gorm:"ForeignKey:UserID" json:"stream"`
}
