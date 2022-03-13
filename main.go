package main

import (
	"log"
	"uniassist/handler"

	"github.com/gin-gonic/gin"
)

var port string = ":5000"

func main(){
	r := gin.Default()
	r.Use(CORSPreflightMiddleware())

	r.GET("/", handler.RootHandler) //webpage
	r.GET("/home", handler.HomeHandler) //data (topPosts + userdata + friendlist)

	r.GET("/login", handler.LoginHandler) //webpage
	r.POST("/loginAuth", handler.LoginAuth) //POST, auth

	r.GET("/register", handler.RegisterHandler)//webpage
	r.POST("/registerAuth", handler.RegisterAuth)//POST, auth

	r.GET("/search", handler.SearchPostHandler)//Search post/s with title
	r.GET("/myPosts", handler.MyPostsHandler) //GET MY POSTS (MUST BE LOGGED IN)
	r.GET("/posts/:idPost", handler.ShowPostHandler) //webpage, gives data needed for displaying a post
	r.GET("/post", handler.PostHandler) //webpage, gives categories for "post form"/"question form"
	r.POST("/postAuth", handler.PostAuth)//POST, auth (title, content, id_category)
	r.DELETE("/posts/:idPost/deleteAuth", handler.DeletePostAuth) //DELETE a post (including the answers)


	r.POST("/posts/:idPost/love",handler.GiveLoveHandler) //give love to a post (backend harus nerima love_value)
	r.DELETE("/posts/:idPost/unlove", handler.UnloveHandler) //unlove a post
	r.PUT("/posts/:idPost/love", handler.UpdateLovePost) // change love value

	r.GET("/posts/:idPost/answer", handler.AnswerHandler) // answer "form", sebenarnya bisa di ShowPostHandler.. tapi kalau mau dipake untuk buat page terpisah gapapa
	r.POST("/posts/:idPost/answerAuth", handler.AnswerAuth) //Answer auth
	r.DELETE("/answer/:idAnswer/deleteAuth", handler.DeleteAnswerAuth) //DELETE an answer

	r.POST("/follow/user/:id", handler.FollowFriend) //Follow a friend
	r.DELETE("/unfollow/user/:id", handler.Unfollow) //Unfollow a friend

	r.PUT("/posts/:idPost/isAnswered", handler.IsAnsweredHandler) //update is answered attribute on a post
	//============-==-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=`-=`-=`-=`-=`-`=-`=-`=-`=-`=-`=-`=-`=`-=`-=`-=`-=`-


	r.GET("/debug", handler.Debug)

	log.Fatal(r.Run(port))

}

func CORSPreflightMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Max-Age", "86400")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, user-info")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        if c.Request.Method == "OPTIONS" {
            c.Writer.Header().Set("Content-Type", "application/json")
            c.AbortWithStatus(204)
        } else {
            c.Next()
        }
    }
}

// //-------------------------migration
// import (
// 	"uniassist/repo"
// )
// func main(){
// 	repo.Migration()
// 	// service.AddCategory()
// }