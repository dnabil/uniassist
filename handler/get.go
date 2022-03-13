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
	loggedIn := false
	tkn, claims, _ :=cookieChecker(c)
	if tkn != nil || claims != nil{
		loggedIn = true
	}
	topPosts, err := service.GetTopPost()
	if err != nil {c.JSON(http.StatusInternalServerError, helper.JsonMessage("ERROR", "Contact administrator"));return;}
	
	if !loggedIn {
		fmt.Println("/n not logged in! /homeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
		c.JSON(http.StatusOK, gin.H{
			"status" : "SUCCESS",
			"topPosts" : topPosts,
		})
		return
	}
	
	log.Println()
	log.Println("UDAH LOGIN YAAAAAAAAAAAAAAAAAAAAAAA")
	log.Println()
	idUser := service.ResponseUserData(claims.Username).ID
	userNFriends := entity.UserNFriends{}
	userdata , friendsdata, err := service.GetUserAndFriendData(idUser)
	if err != nil {c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"));return;}
	userNFriends.User = userdata; userNFriends.Friends = friendsdata;
	
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"userAndFriends" : userNFriends,
		"topPosts" : topPosts,
	})
}


func PostHandler(c *gin.Context){ //post/question form
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil {
		c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))
		return
	}
	
	idUser := service.ResponseUserData(claims.Username).ID
	userNFriends := entity.UserNFriends{}
	userdata , friendsdata, err := service.GetUserAndFriendData(idUser)
	if err != nil {c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))}
	userNFriends.User = userdata; userNFriends.Friends = friendsdata;
	

	//showing categories //JANGAN LUPA GIVE WEBPAGE
	categories := service.GetCategories()
	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"categories" : categories,
		"userAndFriends" : userNFriends,
	})
}


/*
For showing post relative to post's id as parameter
*/
func ShowPostHandler(c *gin.Context){ //show post/question
	loggedIn := false
	tkn, claims, _ := cookieChecker(c)
	if tkn != nil || claims != nil {
		loggedIn = true
	}

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
	//----

	post := entity.Post{}
	post.ID = id
	post ,err = service.GetPost(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.JsonMessage("ERROR", "Post not found"))
		return
	}

	//if user exist
	var idUserLoggedIn uint = 0
	if loggedIn {idUserLoggedIn = service.ResponseUserData(claims.Username).ID}

	//post's
	responseUser, err := service.ResponseUserDataId(post.UserId)
	if err != nil {c.JSON(http.StatusInternalServerError, helper.JsonMessage("ERROR", "Post's owner account may have been deleted"))} 
	responsePost := service.GetResponsePost(&post)
	responseAnswer := service.GetAnswers(post.ID)


	var resp entity.ResponseShowPost
	resp.Post = responsePost
	resp.Category = service.GetCategory(resp.Post.CategoryId)
	resp.User = responseUser
	resp.Answer = responseAnswer
	resp.IsLoved = service.IsLovePost(idUserLoggedIn, responsePost.ID)
	/*rating = (lov / q)*/
	var lov float64=float64(post.Loves); var q float64=float64(service.LoveRowsAffected(responsePost.ID));
	if lov > 0 && q > 0 {resp.Rating = lov / q;} else {resp.Rating = 0}

	fmt.Printf("resp: %v\n", resp)

	if !loggedIn{
		//JANGAN LUPA GIVE WEBPAGE
		c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"post" : resp,
		})
	} else {
		idUser := service.ResponseUserData(claims.Username).ID
		userNFriends := entity.UserNFriends{}
		userdata , friendsdata, err := service.GetUserAndFriendData(idUser)
		if err != nil {c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"));return;}
		userNFriends.User = userdata; userNFriends.Friends = friendsdata;

		c.JSON(http.StatusOK, gin.H{
			"status" : "SUCCESS",
			"post" : resp,
			"userAndFriends" : userNFriends,
			})
	}
}	

//search user by TITLE
func SearchPostHandler(c *gin.Context){
	loggedIn := false
	tkn, claims, _ := cookieChecker(c)
	if tkn != nil || claims != nil {
		loggedIn = true
	}

	q := c.Query("q")
	posts, err := service.SearchPostTitle(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Bad request"))
		return
	}

	if !loggedIn {
		c.JSON(http.StatusOK, gin.H{
			"status" : "SUCCESS",
			"result" : posts,
		})
		return
	}

	idUser := service.ResponseUserData(claims.Username).ID
	userNFriends := entity.UserNFriends{}
	userdata , friendsdata, err := service.GetUserAndFriendData(idUser)
	if err != nil {c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))}
	userNFriends.User = userdata; userNFriends.Friends = friendsdata;

	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"result" : posts,
		"userAndFriends" : userNFriends,
	})
}

//handler for showing user's posts
func MyPostsHandler(c *gin.Context) {
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first.",
		}); log.Println(err.Error())
		return
	}

	myPosts, err := service.GetMyPosts(service.ResponseUserData(claims.Username).ID)
	if err != nil {c.JSON(http.StatusInternalServerError, helper.JsonMessage("ERROR", "Something went wrong. :("));return}

	idUser := service.ResponseUserData(claims.Username).ID
	userNFriends := entity.UserNFriends{}
	userdata , friendsdata, err := service.GetUserAndFriendData(idUser)
	if err != nil {c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"))}
	userNFriends.User = userdata; userNFriends.Friends = friendsdata;

	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"posts" : myPosts,
		"userAndFriends" : userNFriends,
	})
}




func AnswerHandler(c *gin.Context){
	loggedIn := false
	tkn, claims, _ := cookieChecker(c)
	if tkn != nil || claims != nil {
		loggedIn = true
	}

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
	//----

	post := entity.Post{}
	post.ID = id
	post ,err = service.GetPost(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.JsonMessage("ERROR", "Post not found"))
		return
	}

	//if user exist
	var idUserLoggedIn uint = 0
	if loggedIn {idUserLoggedIn = service.ResponseUserData(claims.Username).ID}

	//post's
	responseUser, err := service.ResponseUserDataId(post.UserId)
	if err != nil {c.JSON(http.StatusInternalServerError, helper.JsonMessage("ERROR", "Post's owner account may have been deleted"))} 
	responsePost := service.GetResponsePost(&post)
	responseAnswer := service.GetAnswers(post.ID)


	var resp entity.ResponseShowPost
	resp.Post = responsePost
	resp.Category = service.GetCategory(resp.Post.CategoryId)
	resp.User = responseUser
	resp.Answer = responseAnswer
	resp.IsLoved = service.IsLovePost(idUserLoggedIn, responsePost.ID)
	/*rating = (lov / q)*/
	var lov float64=float64(post.Loves); var q float64=float64(service.LoveRowsAffected(responsePost.ID));
	if lov > 0 && q > 0 {resp.Rating = lov / q;} else {resp.Rating = 0}

	fmt.Printf("resp: %v\n", resp)

	idUser := service.ResponseUserData(claims.Username).ID
	userNFriends := entity.UserNFriends{}
	userdata , friendsdata, err := service.GetUserAndFriendData(idUser)
	if err != nil {c.JSON(http.StatusUnauthorized, helper.JsonMessage("ERROR", "Unauthorized"));return;}
	userNFriends.User = userdata; userNFriends.Friends = friendsdata;

	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"post" : resp,
		"userAndFriends" : userNFriends,
		})

}
