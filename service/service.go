package service

import (
	"fmt"
	"uniassist/entity"
	"uniassist/repo"
)

//AddCategory dipakai jika database kereset/mau nambah kategori
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
}//----

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

func GetPost(id uint) (post entity.Post, err error){
	err = repo.Db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return
	}
	return post, err
}

func DeletePost(post *entity.Post) (err error){
	err = repo.Db.Unscoped().Delete(post).Error
	return
}


/*Get response answers from an id post*/
func GetAnswers(id uint) (answers []entity.ResponseAnswer, err error){
	err = repo.Db.Where("post_id = ?", id).Model(&entity.Answer{}).Find(&answers).Error
	return
}
//Get a single answer entity from AnswerID
func GetAnswerFromId(id uint)(entity.Answer){
	var answer entity.Answer
	err := repo.Db.Where("id = ?", id).Model(&entity.Answer{}).Find(&answer)
	fmt.Printf("err.Error: %v\n", err.Error)
	return answer
}

//delete answers by answer object with idPost filled
func DeleteAnswers(answers *entity.Answer) (err error){
	err = repo.Db.Where("post_id = ?", answers.PostId).Unscoped().Delete(answers).Error
	return
}

//delete a single answer entity
func DeleteAnswer(answer *entity.Answer) (err error) {
	err = repo.Db.Unscoped().Delete(answer).Error
	return
}

//dari username ngasih ResponseUser (ID & USERNAME)
func ResponseUserData(username string) (entity.ResponseUser){
	user := entity.ResponseUser{}
	repo.Db.Where("username = ?", username).Model(&entity.User{}).First(&user)
	return user
}

//dari userId ngasih ResponseUser (ID & USERNAME)
func ResponseUserDataId(id uint) (userData entity.ResponseUser, err error){
	err = repo.Db.Where("id = ?", id).Model(&entity.User{}).First(&userData).Error
	if err != nil {return}
	return userData, nil
}

func FollowFriend(from uint,to uint) (err error){
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

func ShowFriendList(userID uint)(friends []entity.ResponseUser, err error){
	x := []entity.FriendList{}
	err = repo.Db.Where("id = ?", userID).Model(&entity.FriendList{}).Find(&x).Error
	fmt.Printf("x: %v\n", x)
	if err != nil { return }

	FriendID := make([]uint, len(x))
	for i, obj := range x {
		FriendID[i] = obj.FriendID
	}

	err = repo.Db.Where(FriendID).Model(&entity.User{}).Find(&friends).Error
	if err != nil { return }
	return
}

//returns User's stats and user's friends :D (for home's or other pages display)
func GetUserAndFriendData(idUser uint)(userData entity.ResponseUser, friendsData []entity.ResponseUser, err error){
	userData, err = ResponseUserDataId(idUser)
	if err != nil {return}
	friendsData, err = ShowFriendList(idUser)
	return
}


//Fills ResponseShowPost2 struct data
func FillResponseShowPost2(IdPost uint) (result entity.ResponseShowPost2, err error){
	
	// respPost, category, respUser, IsAnswered, Rating
	post, err := GetPost(IdPost)
	if err != nil {return} 
	respPost := GetResponsePost(&post)
	if err != nil {return}
	respUser, err := ResponseUserDataId(post.UserId)
	if err != nil {return}
	category := GetCategory(post.CategoryId)

	result.Post = respPost
	result.User = respUser
	result.Category = category
	result.Rating = 0
	var lov float64=float64(post.Loves); var q float64=float64(LoveRowsAffected(respPost.ID));
	if lov > 0 && q > 0 {result.Rating = lov / q;}

	return
}

//SEARCH POST DATA BY TITLE
func SearchPostTitle(title string) ([]entity.ResponseShowPost2, error){
	title = fmt.Sprint("%" + title +"%")
	type x struct {
		ID uint
	}
	var IdData []x
	err := repo.Db.Where("title LIKE ?", title).Model(&entity.Post{}).Find(&IdData).Error
	if err != nil {return nil, err}
	
	Posts := make([]entity.ResponseShowPost2, len(IdData))
	for i, id := range IdData {
		Posts[i] , err = FillResponseShowPost2(id.ID)
		if err != nil {return Posts, err}
	}
	return Posts, err
}

//GET my Posts
func GetMyPosts(idUser uint) ([]entity.ResponseShowPost2, error){
	type x struct {
		ID uint
	}; var IdData []x;
	err := repo.Db.Where("user_id = ?", idUser).Model(&entity.Post{}).Find(&IdData).Error
	if err != nil {return nil, err}

	Posts := make([]entity.ResponseShowPost2, len(IdData))
	for i, id := range IdData {
		Posts[i], err = FillResponseShowPost2(id.ID)
		if err != nil {return Posts, err}
	}
	return Posts, err
}

//GET LATEST POST DATA
func GetTopPost()([]entity.ResponseShowPost2 ,error){
	var err error = nil
	type x struct {
		ID uint
	}
	var IdData []x
	err = repo.Db.Order("loves desc").Model(&entity.Post{}).Find(&IdData).Error
	if err != nil {return nil, err}

	Posts := make([]entity.ResponseShowPost2, len(IdData))
	 
	for i, id := range IdData {
		Posts[i] , err = FillResponseShowPost2(id.ID)
		if err != nil {return Posts, err}
	}
	return Posts, err
}

//Create a new love id
func CreateLovePost(lovepost entity.LovePost) (err error) {
	err = repo.Db.Create(&lovepost).Error
	if err != nil {
		return err
	}
	return nil
}

//get entity.lovePost{} from db, 
func GetLovePost(userId uint, postId uint) (lovePost entity.LovePost, err error){
	err = repo.Db.Where("user_id = ? AND post_id = ?", userId, postId).First(&lovePost).Error
	if err != nil {
		return
	}
	return lovePost, nil
}

//search if post already loved or not. receives : (UserId, PostId)
func IsLovePost(user uint, post uint) bool{
	lovePost := entity.LovePost{}
	err := repo.Db.Where("user_id = ? AND post_id = ?", user, post).First(&lovePost).Error
	if err == nil {
		return true
	}
	return false
}

//Unlove
func UnLove(lovePost *entity.LovePost) (err error){
	err = repo.Db.Delete(&lovePost).Error
	if err != nil {
		return err
	}
	return nil
}

//UPDATE the loves on post 
func UpdateLovePost(loveValue uint, post *entity.Post, isPlus bool)(err error) {
	past := post.Loves
	if isPlus {
		post.Loves = past + loveValue
	} else {
		post.Loves = past - loveValue
	}
	err = repo.Db.Model(&entity.Post{}).Where("id = ?", post.ID).Update("loves", post.Loves).Error
	if err != nil {
		return
	}
	return nil
}

//UPDATE the love_value on love entity
func UpdateLoveEntity(userId uint, postId uint, love_value uint) (err error){
	err = repo.Db.Model(&entity.LovePost{}).Where("user_id = ? AND post_id = ?", userId, postId).Update("love_value", love_value).Error
	if err != nil {return}
	return nil
}

//Deletes all love from a post, MUST INSERT POST ID IN THE LOVEPOST ENTITY
func DeleteAllLoveFromPost(love *entity.LovePost) (err error) {
	err = repo.Db.Where("post_id = ?", love.PostId).Unscoped().Delete(love).Error
	return
}

//return rows affected for calculating "rating" on post
func LoveRowsAffected(idPost uint) (int64){
	q := repo.Db.Where("post_id = ?", idPost).Find(&[]entity.LovePost{}).RowsAffected
	return q
}
