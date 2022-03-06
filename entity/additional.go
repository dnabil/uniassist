package entity

import "time"

type InputLogin struct {
	Username string `json:"username"` //username or email is allowed
	Password string `json:"password"`
}

//Response
type ResponseLogin struct {
	ID       uint
	Username string
	Email    string
	Name     string

	// Bio Bio
}

type ResponseUser struct {
	ID       uint `json:"id_user"`
	Username string
}

type InputRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InputPost struct {
	Title   string `gorm:"default:'Blank Title Question'" json:"title"` //title question
	Content string `gorm:"default:'Blank Question'" json:"content"`
	UserId  uint   `json:"id_user"`
}

type ResponsePost struct {
	CreatedAt time.Time
	Title string `json:"title"`
	Content    string `json:"content"`

	CategoryId uint `json:"id_category"`
	CategoryName string
}

type ResponseAnswer struct {
	/*search in database first*/
	Username  string //user's
	Content string `json:"answer_content"`
}