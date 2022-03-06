package service

import (
	"log"
	"uniassist/entity"
	"uniassist/repo"
)

func GetCategories() ([]entity.Category){

	var categories []entity.Category
	_ = repo.Db.Find(&categories)
	return categories
}

func GetCategory(id uint)(entity.Category){
	var category entity.Category
	_ = repo.Db.Where("id = ?", id).First(&category)
	return category
}

/*Get response post from an id post*/
func GetPost(post *entity.Post) (responsePost entity.ResponsePost){

	//searching for post
	_ = repo.Db.Where("id = ?", post.ID).Model(&entity.Post{}).First(&responsePost)
	log.Printf("responsePost: %v\n", responsePost)
	log.Printf("post: %v\n", post)
	//username
	var responseUser entity.ResponseUser
	_ = repo.Db.Where("id = ?", post.UserId).Model(&entity.User{}).First(&responseUser)
	return
}


/*Get response answers from an id post*/
func GetAnswers(id uint) ([]entity.ResponseAnswer){
	var answers []entity.ResponseAnswer
	_ = repo.Db.Where("post_id = ?", id).Model(&entity.Answer{}).First(&answers)

	return answers
}

//dari username ngasih ResponseUser (ID & USERNAME)
func ResponseUserData(username string) (entity.ResponseUser){
	user := entity.ResponseUser{}
	repo.Db.Where("username = ?", username).Model(&entity.User{}).First(&user)
	return user
}

//dari userId ngasih ResponseUser (ID & USERNAME)
func ResponseUserDataId(id uint) (entity.ResponseUser){
	user := entity.ResponseUser{}
	repo.Db.Where("id = ?", id).Model(&entity.User{}).First(&user)
	return user
}
