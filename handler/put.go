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

func UpdateLovePost(c *gin.Context) {
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to love a post",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	inputLove := entity.LovePost{} //receives "love_value"
	err = c.ShouldBindJSON(&inputLove)
	var falseLove bool = (inputLove.LoveValue <= 0 || inputLove.LoveValue > 5) 
	if err != nil || falseLove{
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Bad Request"))
		return
	}

	responseUser := service.ResponseUserData(claims.Username)

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
	//----------

	post, err := service.GetPost(id); emptyPost := entity.Post{}
	if post == emptyPost || err != nil{
		c.JSON(http.StatusNotFound, helper.JsonMessage("ERROR", "Post Not Found"))
		return;}

	lovePost, err := service.GetLovePost(responseUser.ID, post.ID)
	if err != nil { c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Bad Request"));
		return;}
	
	err = service.UpdateLoveEntity(lovePost.UserId, lovePost.PostId, inputLove.LoveValue)
	if err != nil {c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Contact administrator")); 
		return;}

	err = service.UpdateLovePost(lovePost.LoveValue, &post, false)
	if err != nil {
		c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Contact administrator"))
		return;}
	err = service.UpdateLovePost(inputLove.LoveValue, &post, true)
	if err != nil {
		c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Contact administrator"))
		return;}

	c.JSON(http.StatusResetContent, helper.JsonMessage("SUCCESS", "Changed love value"))

}