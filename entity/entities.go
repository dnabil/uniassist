package entity

import "time"

type Bio struct {
	ID         uint   `gorm:"primaryKey" json:"id_bio"`
	Content    string `gorm:"default:'The user hasn't wrote anything yet...'" json:"content"`
	Occupation string `gorm:"default:'-'" json:"occupation"`
}

type User struct {
	ID        uint `gorm:"primaryKey" json:"id_user"`
	CreatedAt time.Time
	Username  string
	Email     string
	Password  string
	Name      string `gorm:"default:'user#-'"`

	Bio   Bio 
	BioId uint `json:"id_bio"`
}

type Post struct {
	ID uint `gorm:"primaryKey" json:"id_post"`
	CreatedAt time.Time
	Title    string `gorm:"default:'Blank Title Question'" json:"title"` //title question
	Content    string `gorm:"default:'Blank Question'" json:"content"`

	User User 	
	UserId uint `json:"id_user"`
	
	Category Category
	CategoryId uint `json:"id_category"`
	CategoryName string
}

type Category struct{ //1 post 1 kategori
	ID uint `gorm:"primaryKey" json:"id_category"`
	Name string
}

type Answer struct {
	ID uint `gorm:"primaryKey" json:"id_answer"`
	CreatedAt time.Time
	Content    string `gorm:"default:'Blank Answer'" json:"answer_content"`
	Username string `gorm:"not null"` //received from claims/token

	User User 	
	UserId uint `json:"id_user"`

	Post Post
	PostId uint `json:"id_post"`
}