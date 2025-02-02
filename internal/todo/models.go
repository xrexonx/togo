package todo

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UserId      string `json:"userId"`
}
