package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// "net/url"
	"strconv"
	"time"

	"uniassist/entity"
	"uniassist/helper"
	"uniassist/repo"
	"uniassist/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
var jwtKey = []byte("KeyUniAssist")

//konsistensi
// log.Println(err.Error())
// 	c.JSON(http.StatusUnauthorized, gin.H{
// 		"status":"ERROR"/"SUCCESS",
// 		"message":err.Error(),
// 	})


//tkn, claims :=cookieChecker(c)
func cookieChecker(c *gin.Context) (tkn *jwt.Token, claims *Claims, err error){

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return
		}
		return
	}

	tokenStr := cookie.Value
	claims = &Claims{}

	tkn, err = jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return
		}
		return
	}
	if !tkn.Valid {
		return
	}
	return //tkn, claims and err
}


func RegisterAuth(c *gin.Context){
	input := entity.InputRegister{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "error, please contact administrator",
		})
		log.Println("register json binding failed")
		return
	}

	theUser := entity.User{}

	//checking blank
	isBlank := func(data string) (err bool){
		if data == "" {
			err = true
			return
		}
		err = false
		return
	}
	if isBlank(input.Email) {c.JSON(http.StatusBadRequest, gin.H{"error" : "Do not leave Email blank"}); return}
	if isBlank(input.Password) {c.JSON(http.StatusBadRequest, gin.H{"error" : "Do not leave Password blank"}); return}
	if isBlank(input.Username) {c.JSON(http.StatusBadRequest, gin.H{"error" : "Do not leave Username blank"}); return}

	/*email and username must be different from the ones that exist in database*/
	//checking email
	repo.Db.Where("email = ?", input.Email).First(&theUser)
	if theUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Email already used",
		})
		return
	}

	//checking username
	repo.Db.Where("username = ?", input.Username).First(&theUser)
	if theUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Username already used",
		})
		return
	}
	//success, creating new account/user
	theUser.Name = fmt.Sprintf("user#%d", theUser.ID)
	theUser.Email = input.Email
	theUser.Username = input.Username
	//storing hashed password
	theUser.Password, err = HashNSalt([]byte(input.Password)) //we don't store your "password" here :)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.JsonMessage("ERROR", "Contact administator"))
	}
	repo.Db.Create(&theUser)
	
	c.JSON(http.StatusOK, gin.H{
		"message" : "Registered as",
		"account" : input,
	})
}

func LoginAuth(c *gin.Context){
	var InputLogin  entity.InputLogin
	var User entity.User

	err := c.ShouldBindJSON(&InputLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message":err.Error(),
		})
		log.Println("JSON binding failed")
		return
	}

	//finding username
	_ = repo.Db.Where("username = ?", InputLogin.Username).First(&User).Error
	var isExist = false
	if User.Username != ""{
		isExist = true
	}
	if !isExist{
		//finding email
		repo.Db.Where("email = ?", InputLogin.Username).First(&User)
		var isExist = false
		if User.Email != ""{
		isExist = true
		}
		if !isExist{
			c.JSON(http.StatusBadRequest, gin.H{
				"status":"ERROR",
				"message":"Invalid email/username",
			})
			return
		}
	}

	//comparing password
	passError := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(InputLogin.Password)) 
	if passError == bcrypt.ErrMismatchedHashAndPassword && passError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message":"Invalid password",
		})
		return
	}

	/*TOKEN's DEMISE >:D*/
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: InputLogin.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":"ERROR",
			"message":"Error signing token, contact administrator",
		}); log.Println(err.Error())
		return
	}

	http.SetCookie(c.Writer,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		
	c.JSON(http.StatusNotFound, gin.H{
		"status":"SUCCESS",
		"message" : "Login success",
	});
}


func HashNSalt(pass []byte,) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost) //(args,Hashmethod)
	if err != nil {
		log.Fatalln("failed hashing password")
	}

	return string(hashed), err
}


/*
Post authorization (MUST BE LOGGED IN)
*/
func PostAuth(c *gin.Context) {

	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{ //if not authorized, redirect to login page
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to post a question",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	responseUser := service.ResponseUserData(claims.Username)
	log.Println("===============for post, RESPONSE USER :", responseUser) //debug

	post := entity.Post{}
	err = c.ShouldBindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message":"Bad request",
		}); log.Println(err.Error())
		return
	}
	// fmt.Printf("post.Title: %v\n", post.Title) //debug
	// fmt.Printf("post.Content: %v\n", post.Content) //debug
	if post.Title == "" {c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message" : "Title shouldn't be empty",
		}); 
		return}
	if post.Content ==  "" {c.JSON(http.StatusBadRequest, gin.H{
		"status":"ERROR",
		"message" : "Content shouldn't be empty",
	}); 
	return}
	
	post.UserId = responseUser.ID

	err = repo.Db.Create(&post).Error
	if err != nil { log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message":"Failed creating object post",
		})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"status" : "SUCCESS",
		"message" : "new post created",
		"data" : post,
	})
}


func AnswerAuth(c *gin.Context) {
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to answer a question/post",
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
	//end of convert
	//searching for post
	post := entity.Post{}
	post.ID = id
	errdb := repo.Db.Where("id = ?", post.ID).First(&post)
	if errdb.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":"ERROR",
			"message" : "404 Post not found",
		}); log.Println(err.Error())
		return
	} //--------


	answer := entity.Answer{} //receives "content" and username
	answer.Username = claims.Username
	err = c.ShouldBindJSON(&answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message" : "bad request",
		}); log.Println(err.Error())
		return
	}
	if answer.Content == ""{
		c.JSON(http.StatusBadRequest, gin.H{
			"status":"ERROR",
			"message" : "Answer shouldn't be empty",
		}); log.Println(err.Error())
		return
	}

	responseUser := service.ResponseUserData(claims.Username)
	answer.UserId = responseUser.ID
	answer.PostId = id

	errdb = repo.Db.Create(&answer)
	if errdb.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":"ERROR",
			"message" : "Failed to create object answer, contact administrator",
		}); log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status" : "SUCCESS",
		"message" : "Answer created",
	})
}

//search user by username
func SearchPostHandler(c *gin.Context){
	q := c.Query("q")
	q = strings.ReplaceAll(q, " ", "%")
	posts := service.SearchPostTitle(q)
	if posts == nil {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "No posts found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : "SUCCESS",
		"message" : "Data found",
		"result" : posts, 
	})
}

//follow friend by id
func FollowFriend(c *gin.Context) {
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

	//if USER tried to FOLLOW HIMSELF, reject the narcissistic bastard, sorry forgive me
	if user.ID == id {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "narcissism"))
		return
	}
	//--------------

	friend := service.ResponseUserDataId(id)
	var isBlank bool= (friend == entity.ResponseUser{}) 
	if isBlank {
		c.JSON(http.StatusNotFound, helper.JsonMessage("ERROR", "ID NOT FOUND"))
		return
	}

	err = service.FollowFriend(user.ID, id)
	if err != nil {
		c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Already followed. if not, contact the administrator"))
		return
	}

	c.JSON(http.StatusOK, helper.JsonMessage("SUCCESS", "Followed :D"))
}



//catatan kecil 

//CARA REDIRECT
// https://stackoverflow.com/questions/61970551/golang-gin-redirect-and-render-a-template-with-new-variables
// // first solution
// c.SetCookie("wage", "123", 10, "/", c.Request.URL.Hostname(), false, true)
// c.SetCookie("amount", "13123", 10, "/", c.Request.URL.Hostname(), false, true)
// location := url.URL{Path: "/api/callback/cookies",}
// c.Redirect(http.StatusFound, location.RequestURI())

// // second solution
// q := url.Values{}
// q.Set("wage", "123")
// q.Set("amount", "13123")
// location := url.URL{Path: "/api/callback/query_params", RawQuery: q.Encode()}
// c.Redirect(http.StatusFound, location.RequestURI())