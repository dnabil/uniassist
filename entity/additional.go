package entity

import (
	"gorm.io/gorm"
)

type InputLogin struct {
	Email string `json:"email"`
	Username string `json:"username"`
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
	IsAnswered bool `gorm:"default:false"`

	CategoryId uint
}

type ResponseShowPost struct {
	Post ResponsePost
	Category Category
	User ResponseUser
	Answer []ResponseAnswer

	//tambahan, diisi pake method di service
	IsLoved bool
	Rating float64
}

type ResponseShowPost2 struct {
	Post ResponsePost
	Category Category
	User ResponseUser

	//tambahan, diisi pake method di service
	Rating float64
}

type ResponseAnswer struct {
	/*search in database first*/
	gorm.Model
	Username  string //user's
	Name string //user's
	Content string
}
type ResponseLovePost struct {
	ID uint
	LoveValue uint
	PostId uint
	UserId uint
}

type UserNFriends struct {
	User ResponseUser
	Friends []ResponseUser
}