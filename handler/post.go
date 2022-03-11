package handler

import (
	"fmt"
	"log"
	"net/http"

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
func cookieChecker(c *gin.Context) (*jwt.Token, *Claims, error){

	claims2 := &Claims{}
	tkn2 :=  &jwt.Token{}
	tokenStr2 := c.Request.Header.Get("token")
	// println()
	// println("str2 : ", tokenStr2)
	tkn2, err := jwt.ParseWithClaims(tokenStr2, claims2,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return nil, nil, fmt.Errorf("Unauthorized")
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return nil, nil, fmt.Errorf("Bad request")
	}

	if !tkn2.Valid {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return nil, nil, fmt.Errorf("Unauthorized")
	}
		
	println(claims2.Username)
	return tkn2, claims2, err
}

func Debug(c *gin.Context){
	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{ //if not authorized, redirect to login page
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to post a question",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, "WELCOME")
}


func RegisterAuth(c *gin.Context){
	input := entity.InputRegister{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Contact administrator"))
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
	theUser.Name = fmt.Sprintf("%s %s",input.Username , "Subagyo")
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
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Contact Administrator"))
		log.Println("JSON binding failed")
		return
	}
	
	//finding username
	err = repo.Db.Where("username = ?", InputLogin.Username).First(&User).Error
	if err != nil && User.Username == "" {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Username shouldn't be empty")) ;return;}
	if err != nil {c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "wrong username/email"));return;}
	
	//finding email
	err = repo.Db.Where("email = ?", InputLogin.Email).First(&User).Error
	if err != nil && User.Email == ""{
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Email shouldn't be empty"));return;}
	if err != nil {c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "wrong username/email"));return;}
	

	//comparing password
	passError := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(InputLogin.Password)) 
	if passError == bcrypt.ErrMismatchedHashAndPassword && passError != nil {
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Wrong Password, try again"))
		return
	}

	/*TOKEN's DEMISE >:D*/
	expirationTime := time.Now().Add(time.Minute * 30)

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
		
	c.JSON(http.StatusOK, gin.H{
		"status":"SUCCESS",
		"message" : "Login success",
		"token" : tokenString,
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
		c.JSON(http.StatusNotImplemented, gin.H{
			"status":"ERROR",
			"message":"Failed creating object post",
		})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"status" : "SUCCESS",
		"message" : "New post created",
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
	answer.Name = responseUser.Name
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

	friend, _ := service.ResponseUserDataId(id)
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

func GiveLoveHandler(c *gin.Context){

	tkn, claims, err := cookieChecker(c)
	if tkn == nil || claims == nil || err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":"ERROR",
			"message" : "Please register/login first to give love to a post",
			"data" : fmt.Sprintf("%s%s",c.Request.Host,"/login"),
		}); log.Println(err.Error())
		return
	}

	inputLove := entity.LovePost{}
	err = c.ShouldBindJSON(&inputLove)
	var falseLove bool = (inputLove.LoveValue <= 0 || inputLove.LoveValue > 5) 
	if err != nil || falseLove{
		c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "Bad Request"))
		return
	}

	responseUser := service.ResponseUserData(claims.Username)
	userid := responseUser.ID

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

	post, err := service.GetPost(id); emptyPost := entity.Post{}
	fmt.Printf("id: %v\n", id)
	fmt.Printf("post: %v\n", post)
	if post == emptyPost {
		c.JSON(http.StatusNotFound, helper.JsonMessage("ERROR", "Post Not Found"))
		return;}

	loveD := service.IsLovePost(userid, id)	
	if loveD {c.JSON(http.StatusBadRequest, helper.JsonMessage("ERROR", "You can't do that")); return}

	inputLove.PostId = id
	inputLove.UserId = userid

	err = service.CreateLovePost(inputLove)
	if err != nil {
		c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Contact administrator"))
		return;}
	err = service.UpdateLovePost(inputLove.LoveValue, &post, true);
	if err != nil {
		c.JSON(http.StatusNotImplemented, helper.JsonMessage("ERROR", "Contact administrator"))
		return;}
	
	fmt.Printf("%s give %v love to post id %v\n",claims.Username , inputLove.LoveValue, inputLove.PostId)
	c.JSON(http.StatusOK, helper.JsonMessage("SUCCESS", "Love given"))
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