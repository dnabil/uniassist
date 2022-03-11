package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"not null" type:"VARCHAR(255)"`
	Email     string `gorm:"not null" type:"VARCHAR(255)"`
	Password  string `gorm:"not null"`
	Name      string `gorm:"default:'user#-'"`
	// Users []User `gorm:"many2many:friendlist;"`
	
	ProfilePic string `gorm:"default:'https://freepikpsd.com/file/2019/10/default-profile-image-png-1-Transparent-Images.png'"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"not null" json:"title"` //title question
	Content    string `gorm:"not null" json:"content"`

	User User 	
	UserId uint `gorm:"not null"`
	
	Category Category
	CategoryId uint `gorm:"default:1" json:"category_id"`

	Loves uint `gorm:"default:0"`//from LovePost (akumulasi)
	IsAnswered bool `gorm:"default:false"`
}

type Category struct{ //1 post 1 kategori (ngikutin dari desain)
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
type Answer struct {
	gorm.Model
	Content    string `gorm:"not null"`//`gorm:"default:'Blank Answer'"`
	Username string `gorm:"not null"` //received from claims/token
	Name string `gorm:"not null"` //received from claims/token

	User User 	
	UserId uint `gorm:"not null"`

	Post Post
	PostId uint `gorm:"not null"`
}

type LovePost struct{
	ID        uint `gorm:"primarykey"`
	LoveValue uint `json:"love_value"`

	Post Post 
	PostId uint 

	User User 
	UserId uint 
}