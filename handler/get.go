package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"uniassist/entity"
	"uniassist/repo"
	"uniassist/service"

	"github.com/gin-gonic/gin"
)

/*
Jika belum login, maka tampilkan landing page atau "/"
Jika SUDAH LOGIN, redirect ke "/home"
*/
func Root(c *gin.Context) {
	tkn, claims :=cookieChecker(c)
	if tkn == nil || claims == nil { 
		c.JSON(http.StatusOK, gin.H{
			"title":   "UniAssist, Belajar lebih asik!",
			"content": "AYO JOIN KAMI SEKARANG",
		})
		return
	}

	location := url.URL{Path: "/home",}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func Login(c *gin.Context) { // /loginAuth
	tkn, claims := cookieChecker(c)
	if tkn != nil && claims != nil{ //if cookie exist, redirect to home
		location := url.URL{Path: "/home",}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}

	//MASUKKAN PAGE DISINI
	c.JSON(http.StatusOK, gin.H{
		"message" : "welcome to login page:D",
	})
}

func Register(c *gin.Context) { // /registerAuth
	tkn, claims := cookieChecker(c)
	if tkn != nil && claims != nil{ //if cookie exist, redirect to home
		location := url.URL{Path: "/home",}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}

	//MASUKKAN PAGE DISINI
	c.JSON(http.StatusOK, gin.H{
		"message" : "welcome to register page:D",
	})
}

func Home(c *gin.Context) { //home handler
	tkn, claims :=cookieChecker(c)
	if tkn == nil || claims == nil{
		location := url.URL{Path: "/login",}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}

	c.Writer.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))
	fmt.Printf("tkn: %v\n", tkn)
	fmt.Printf("claims: %v\n", claims)
}


func Post(c *gin.Context){ //post/question form
	tkn, claims := cookieChecker(c)
	if tkn == nil || claims == nil{
		//do something atau PAKAI TOKEN NAMA LAIN BIAR REGISTERED USER BISA JAWAB
		return
	}
	
	//showing categories
	categories := service.GetCategories()
	c.JSON(http.StatusOK, categories)
}


/*
For showing post relative to post's id as parameter
*/
func ShowPost(c *gin.Context){ //show post/question
	idPost := c.Param("idPost")
	idPostInt, err := strconv.Atoi(idPost)
	fmt.Printf("err: %v\n", err)
	if err != nil || idPostInt < 0{
		c.JSON(http.StatusNotFound, "404 Not found")
		return
	}; var id uint = uint(idPostInt)

	post := entity.Post{}
	post.ID = id
	errdb := repo.Db.Where("id = ?", post.ID).First(&post)
	if errdb.Error != nil {
		c.JSON(http.StatusNotFound, "Post/Question Not found")
		return
	}
	
	responseUser := service.ResponseUserDataId(post.UserId)
	responsePost := service.GetPost(&post)
	responseAnswer := service.GetAnswers(post.ID)
	answeButton := fmt.Sprintf("%s%s%s",c.Request.Host ,c.Request.URL.Path, "/answer")


	c.JSON(http.StatusOK, gin.H{
		"user" : responseUser,
		"post" : responsePost,
		"answers" : responseAnswer,
		"answerButton" : answeButton,
	})
}	


func Answer(c *gin.Context){
	tkn, claims := cookieChecker(c)
	if tkn == nil || claims == nil{
		location := url.URL{Path: "/login",}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}


	//converting idPost to uint
	idPost := c.Param("idPost")
	idPostInt, err := strconv.Atoi(idPost)
	fmt.Printf("err: %v\n", err)
	if err != nil || idPostInt < 0{
		c.JSON(http.StatusNotFound, "404 Not found")
		return
	}; var id uint = uint(idPostInt) //id == idPost
	//end of convert

	//searching for post
	post := entity.Post{}
	post.ID = id
	errdb := repo.Db.Where("id = ?", post.ID).First(&post)
	if errdb.Error != nil {
		c.JSON(http.StatusNotFound, "Post/Question Not found")
		return
	} //--------
	
	responseUser := service.ResponseUserDataId(post.UserId)
	responsePost := service.GetPost(&post)

	//JANGAN LUPA GIVE WEBPAGE
	c.JSON(http.StatusOK, gin.H{
		"user" : responseUser,
		"post" : responsePost,
	})

}