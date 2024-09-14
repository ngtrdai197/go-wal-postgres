package model

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title            string `json:"title" gorm:"title"`
	ShortDescription string `json:"short_description" gorm:"short_description"`
	Body             string `json:"body" gorm:"body"`
	// ...
}
