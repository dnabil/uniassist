package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"uniassist/entity"
	"uniassist/helper"
	"uniassist/repo"
	"uniassist/service"

	"github.com/gin-gonic/gin"
)

/*
Jika belum login, maka tampilkan landing page atau "/"
Jika SUDAH LOGIN, redirect ke "/home"
*/
func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title":   "UniAssist, Belajar lebih asik!",
		"content": "AYO JOIN KAMI SEKARANG",
	})
}

func LoginHandler(c *gin.Context) {
	tkn, claims, _ := cookieChecker(c)
	if tkn != nil && claims != nil {
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))
		return
	}
	
	//JANGAN LUPA GIVE WEBPAGE
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"data" : fmt.Sprintf("%s%s", c.Request.Host,"/loginAuth"),
	})
}

func RegisterHandler(c *gin.Context) { // /registerAuth
	tkn, claims, _ := cookieChecker(c)
	if tkn != nil && claims != nil{
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))
		return
	}

	//JANGAN LUPA GIVE WEBPAGE
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"data" : fmt.Sprintf("%s%s", c.Request.Host,"/registerAuth"),
	})
}

func HomeHandler(c *gin.Context) { //home handler
	tkn, claims, _ :=cookieChecker(c)
	if tkn == nil || claims == nil{
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
	})
}


func PostHandler(c *gin.Context){ //post/question form
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil {
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))
		return
	}
	
	//showing categories //JANGAN LUPA GIVE WEBPAGE
	categories := service.GetCategories()
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"categories" : categories,
	})
}


/*
For showing post relative to post's id as parameter
*/
func ShowPostHandler(c *gin.Context){ //show post/question
	idPost := c.Param("idPost")
	idPostInt, err := strconv.Atoi(idPost)
	fmt.Printf("err: %v\n", err)
	if err != nil || idPostInt < 0{
		c.JSON(http.StatusNotFound, gin.H{
			"status" : "ERROR",
			"message" : " 404 Post NOT FOUND",
		})
		return
	}; var id uint = uint(idPostInt) // id==idPost

	post := entity.Post{}
	post.ID = id
	errdb := repo.Db.Where("id = ?", post.ID).First(&post)
	if errdb.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status" : "ERROR",
			"message" : "404 Post NOT FOUND",
		})
		return
	}

	responseUser := service.ResponseUserDataId(post.UserId)
	responsePost := service.GetResponsePost(&post)
	responseAnswer := service.GetAnswers(post.ID)

	fmt.Printf("responseUser: %v\n", responseUser)
	fmt.Println()
	fmt.Printf("responsePost: %v\n", responsePost)
	fmt.Println()
	fmt.Printf("responseAnswer: %v\n", responseAnswer)
	fmt.Println()
	// responseUser := service.ResponseUserDataId(post.UserId)
	// responsePost := service.GetResponsePost(&post)
	// responseAnswer := service.GetAnswers(post.ID)
	answerButton := fmt.Sprintf("%s%s%s",c.Request.Host ,c.Request.URL.Path, "/answer")

	var resp entity.ResponseShowPost
	resp.Post = responsePost
	resp.User = responseUser

	fmt.Printf("resp: %v\n", resp)

	//JANGAN LUPA GIVE WEBPAGE
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"post" : resp,
		"answers" : responseAnswer,
		"answerButton" : answerButton,
	})
}	


func AnswerHandler(c *gin.Context){
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))
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
	responsePost := service.GetResponsePost(&post)
	var resp entity.ResponseShowPost
	resp.Post = responsePost
	resp.User = responseUser

	//JANGAN LUPA GIVE WEBPAGE
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESSS",
		"post" : resp,
	})

}
