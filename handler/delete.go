package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"uniassist/entity"
	"uniassist/helper"
	"uniassist/service"

	"github.com/gin-gonic/gin"
)

func DeletePostAuth(c *gin.Context) {
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{ //if not authorized, redirect to login page
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to post a question",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	//converting idPost to uint
	idPost := c.Param("idPost")
	idPostInt, err := strconv.Atoi(idPost)

	if err != nil || idPostInt < 0{
		c.JSON(http.StatusNotFound, gin.H{
			"status":"ERROR",
			"message" : "404 not found",
		}); log.Println(err.Error())
		return
	}; var id uint = uint(idPostInt) //id == idPost

	post := service.GetPost(id);
	if post.ID <= 0 {
		c.JSON(http.StatusNotFound, helper.JsonMessage("ERROR", "Post not found"))
		return
	}
	responseUser := service.ResponseUserData(claims.Username)

	// log.Println(responseUser.ID)
	// log.Println(post.UserId)
	if responseUser.ID != post.UserId {
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "unauthorized"))
		return
	}
	answers := entity.Answer{PostId: post.ID}
	service.DeleteAnswers(&answers)
	service.DeletePost(&post)
}

func DeleteAnswerAuth(c *gin.Context) {
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{ //if not authorized, redirect to login page
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to post a question",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	//converting idPost to uint
	idAnswer := c.Param("idAnswer")
	idAnswerInt, err := strconv.Atoi(idAnswer)

	if err != nil || idAnswerInt < 0{
		c.JSON(http.StatusNotFound, gin.H{
			"status":"ERROR",
			"message" : "404 not found",
		}); log.Println(err.Error())
		return
	}; var id uint = uint(idAnswerInt) //id == idAnswer

	answer := service.GetAnswerFromId(id)
	fmt.Printf("answer: %v\n", answer)
	fmt.Printf("claims.Username: %v\n", claims.Username)
	if answer.Username != claims.Username {
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "unauthorized"))
		return
	} 

	service.DeleteAnswer(&answer);
}

func Unfollow(c *gin.Context){
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to follow someone.",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	//converting param to uint
	idUser := c.Param("id")
	idUserInt, err := strconv.Atoi(idUser)

	if err != nil || idUserInt < 0{
		c.JSON(http.StatusNotFound, gin.H{
			"status":"ERROR",
			"message" : "404 not found",
		}); log.Println(err.Error())
		return
	}; var id uint = uint(idUserInt) //id == "friend"'s id
	//end of convert


	user := service.ResponseUserData(claims.Username)
	//if USER tried to UNFOLLOW HIMSELF, return something nice
	if user.ID == id {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Love yourself:D"))
		return
	}
	//--------------
	err = service.UnfollowFriend(user.ID, id)
	if err != nil {
		c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Contact the administrator"))
		return
	}

	c.JSON(http.StatusOK, helper.JsonMessage("SUCCESS", "Unfollowed"))
}