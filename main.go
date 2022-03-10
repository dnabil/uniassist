package main

import (
	"uniassist/handler"

	"github.com/gin-gonic/gin"
)

var port string = ":8443"

func main(){
	r := gin.Default()

	r.GET("/", handler.RootHandler) //webpage
	r.GET("/home", handler.HomeHandler) //webpage
	
	r.GET("/login", handler.LoginHandler) //webpage
	r.POST("/loginAuth", handler.LoginAuth) //POST, auth
	
	r.GET("/register", handler.RegisterHandler)//webpage
	r.POST("/registerAuth", handler.RegisterAuth)//POST, auth
	
	r.GET("/search", handler.SearchPostHandler)//Search post/s with title
	
	r.GET("/posts/:idPost", handler.ShowPostHandler) //webpage, gives data needed for displaying a post
	r.GET("/post", handler.PostHandler) //webpage, gives categories for "post form"
	r.POST("/postAuth", handler.PostAuth)//POST, auth (title, content, id_category)
	r.DELETE("/posts/:idPost/deleteAuth", handler.DeletePostAuth) //DELETE a post (including the answers)
	
	r.GET("/posts/:idPost/answer", handler.AnswerHandler) //GET DATA, (Answer form) gives data to be displayed
	r.POST("/posts/:idPost/answerAuth", handler.AnswerAuth) //Answer auth
	r.DELETE("/answer/:idAnswer/deleteAuth", handler.DeleteAnswerAuth) //DELETE an answer

	r.POST("/follow/user/:id", handler.FollowFriend) //Follow a friend
	r.DELETE("/unfollow/user/:id", handler.Unfollow) //Unfollow a friend

	r.Run(port)



}

// //-------------------------migration
// import (
// 	"uniassist/service"
// )
// func main(){
// 	// repo.Migration()
// 	service.AddCategory()
// }