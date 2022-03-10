package entity

import (
	"gorm.io/gorm"
)

type InputLogin struct {
	Username string `json:"username"` //username or email is allowed
	Password string `json:"password"`
}

type ResponseUser struct {
	ID       uint 
	Username string
	Name string
	ProfilePic string //a link
}

type InputRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FriendList struct {
	ID uint //user
	FriendID uint //tujuan follow
}

type InputPost struct {
	Title   string `gorm:"default:'Blank Title Question'" json:"title"` //title question
	Content string `gorm:"default:'Blank Question'" json:"content"`
	UserId  uint
	CategoryId uint `gorm:"default:1"`
}

type ResponsePost struct {
	gorm.Model
	Title string `json:"title"`
	Content    string `json:"content"`

	CategoryId uint
}

type ResponseShowPost struct {
	Post ResponsePost
	User ResponseUser
}

type ResponseAnswer struct {
	/*search in database first*/
	gorm.Model
	Username  string //user's
	Content string
}