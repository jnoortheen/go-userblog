package models_test

import "muserblog/models"

var (
	postTitle = "Post 1"
	postContent = "Post 1 content"
)

func postForTest() *models.Post {
	return &models.Post{Title: postTitle, Content: postContent}
}

func (as *ModelSuite) Test_Post() {
	prevCount := as.countObjects(models.Post{})
	post := postForTest()
	as.NoError(as.DB.Create(post))
	as.Equal(as.countObjects(models.Post{}) - prevCount, 1)
}
