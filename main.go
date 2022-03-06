package main

import (
	"uniassist/handler"

	"github.com/gin-gonic/gin"
)

var port string = ":8443"

func main(){
	r := gin.Default()

	r.GET("/", handler.Root) //SUDAH tinggal webpage
	r.GET("/home", handler.Home) //SUDAH tinggal webpage

	r.GET("/register", handler.Register) //SUDAH tinggal webpage
	r.POST("/registerAuth", handler.RegisterAuth)

	r.GET("/login", handler.Login) //SUDAH tinggal webpage
	r.POST("/loginAuth", handler.LoginAuth)

	r.GET("/post", handler.Post) //receiving categories //SUDAH tinggal webpage
	r.POST("/postAuth", handler.PostAuth)
	r.GET("/posts/:idPost", handler.ShowPost) //showing post/question (with the answer/s) with id as parameter //SUDAH tinggal webpage

	r.GET("/posts/:idPost/answer", handler.Answer) //dari button //SUDAH tinggal webpage
	r.POST("/posts/:idPost/answerAuth", handler.AnswerAuth)

	r.Run(port)

}

// // //-------------------------migration
// import (
// 	"uniassist/repo"
// )
// func main(){
// 	repo.Migration()
// }