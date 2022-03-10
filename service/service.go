package service

import (
	"fmt"
	"uniassist/entity"
	"uniassist/repo"
)

func AddCategory() {
	var categories = []string {
	"Other",
	"Pemrograman Lanjut",
	"Pemrograman Dasar",
	"Sistem Basis Data",
	"Jaringan",
	"Pemrograman Web",
	"Pemrograman Mobile",
	"Java",
	"C++",
	"Javascript"}
	
	for i := 0; i < len(categories); i++ {
		obj := entity.Category{}
		obj.Name = categories[i]
		repo.Db.Create(&obj)
	}

}

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

/*Get response post from post entity (with id fiiled)*/
func GetResponsePost(post *entity.Post) (responsePost entity.ResponsePost){

	//searching for response post
	_ = repo.Db.Where("id = ?", post.ID).Model(&entity.Post{}).First(&responsePost)
	// log.Printf("responsePost: %v\n", responsePost)
	// log.Printf("post: %v\n", post)

	//username
	var responseUser entity.ResponseUser
	_ = repo.Db.Where("id = ?", post.UserId).Model(&entity.User{}).First(&responseUser)
	return
}

func GetPost(id uint) (post entity.Post){

	_ = repo.Db.Where("id = ?", id).First(&post)

	return post
}

func DeletePost(post *entity.Post){
	_ = repo.Db.Delete(post).Error
}


/*Get response answers from an id post*/
func GetAnswers(id uint) ([]entity.ResponseAnswer){
	var answers []entity.ResponseAnswer
	_ = repo.Db.Where("post_id = ?", id).Model(&entity.Answer{}).Find(&answers)

	return answers
}
//Get a single answer entity from AnswerID
func GetAnswerFromId(id uint)(entity.Answer){
	var answer entity.Answer
	err := repo.Db.Where("id = ?", id).Model(&entity.Answer{}).Find(&answer)
	fmt.Printf("err.Error: %v\n", err.Error)
	return answer
}

//delete answers by answer object with idPost filled
func DeleteAnswers(answers *entity.Answer){
	_ = repo.Db.Where("post_id = ?", answers.PostId).Delete(answers)
}

//delete a single answer entity
func DeleteAnswer(answer *entity.Answer) {
	_ = repo.Db.Delete(answer)
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

func FollowFriend (from uint,to uint) (err error){
	//filter if already followed, then return redundant
	FriendList := entity.FriendList{}
	err = repo.Db.Where("id = ? AND friend_id = ?", from, to).First(&FriendList).Error
	if err == nil {
		err = fmt.Errorf("Redundant")
		return
	}
	//---------
	FriendList = entity.FriendList{ID: from, FriendID:to}
	err = repo.Db.Create(&FriendList).Error
	return err
}
func UnfollowFriend(from uint,to uint) (err error){
	//filter if not followed, then return redundant
	FriendList := entity.FriendList{}
	err = repo.Db.Where("id = ? AND friend_id = ?", from, to).First(&FriendList).Error
	if err != nil {
		err = fmt.Errorf("Redundant")
		return
	}
	//---------

	err = repo.Db.Delete(&FriendList).Error
	return
}


//SEARCH POST DATA BY TITLE
func SearchPostTitle(title string) ([]entity.Post){
	title = fmt.Sprint("%" + title +"%")
	posts := []entity.Post{}
	repo.Db.Where("title LIKE ?", title).Find(&posts)
	return posts
}

