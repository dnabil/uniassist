package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"uniassist/entity"
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


//tkn, claims :=cookieChecker(c)
func cookieChecker(c *gin.Context) (tkn *jwt.Token, claims *Claims){
	cookie, err := c.Request.Cookie("token")
	fmt.Printf("cookie: %v\n", cookie)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		if err == http.ErrNoCookie {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
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
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	return //token and claims
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
	theBio := entity.Bio{}
	repo.Db.Create(&theBio)
	theUser.Name = fmt.Sprintf("user#%d", theBio.ID)
	theUser.Email = input.Email
	theUser.Username = input.Username
	theUser.Bio = theBio
	//storing hashed password
	theUser.Password = HashNSalt([]byte(input.Password)) //we don't store your "password" here :)
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
			"error" : "Login error, contact the administrator",
		})
		log.Println("JSON binding failed")
		return
	}

	//finding username
	repo.Db.Where("username = ?", InputLogin.Username).First(&User)
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
				"error" : "Invalid email/username",
			})
			return
		}
	}

	//comparing password
	passError := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(InputLogin.Password)) 
	if passError == bcrypt.ErrMismatchedHashAndPassword && passError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Invalid password",
		})
		return
	}

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
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
}

func HashNSalt(pass []byte,) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost) //(args,Hashmethod)
	if err != nil {
		log.Fatalln("failed hashing password")
	}

	return string(hashed)
}

/*
Post authorization (MUST BE LOGGED IN)
*/
func PostAuth(c *gin.Context) {
	tkn, claims := cookieChecker(c)
	if tkn == nil || claims == nil{ //if not authorized, redirect to login page
		location := url.URL{Path: "/login",}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	}
	responseUser := service.ResponseUserData(claims.Username)
	log.Println("===============RESPONSE USER :", responseUser)

	post := entity.Post{}
	c.ShouldBindJSON(&post)
	// fmt.Printf("post.Title: %v\n", post.Title) //debug
	// fmt.Printf("post.Content: %v\n", post.Content) //debug
	if post.Title == "" {c.JSON(http.StatusBadRequest, gin.H{"error" : "title shouldn't be empty"}); return}
	if post.Content == "" {c.JSON(http.StatusBadRequest, gin.H{"error" : "content shouldn't be empty"}); return}
	
	category := service.GetCategory(post.CategoryId)
	post.User.ID = responseUser.ID
	post.User.Username = claims.Username
	post.CategoryName = category.Name
	repo.Db.Create(&post)

	//debug
	c.JSON(http.StatusCreated, gin.H{
		"success" : post,
	})
}

func AnswerAuth(c *gin.Context) {
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


	answer := entity.Answer{} //receives "content" and username
	answer.Username = claims.Username
	err = c.ShouldBindJSON(&answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "error binding json answer",})
		return
	}
	if answer.Content == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error" : "answer shouldn't be empty",})
		return
	}

	responseUser := service.ResponseUserData(claims.Username)
	answer.UserId = responseUser.ID
	answer.PostId = id

	errdb = repo.Db.Create(&answer)
	if errdb.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to create answer, please contact the administator,",})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success" : "answer created",
	})
}


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