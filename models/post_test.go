package models_test

import "muserblog/models"

var (
	postTitle   = "Post 1"
	postContent = "Post 1 content"
)

func postForTest() *models.Post {
	return &models.Post{Title: postTitle, Content: postContent}
}

func countPosts(as *ModelSuite) int {
	postsCount, err := as.DB.Count(models.Posts{})
	as.NoError(err)
	return postsCount
}

func (as *ModelSuite) Test_Post() {
	prevCount := countPosts(as)
	post := postForTest()
	as.NoError(as.DB.Create(post))
	as.Equal(countPosts(as)-prevCount, 1)
}
