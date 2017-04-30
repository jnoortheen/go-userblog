package models_test

import (
	"github.com/satori/go.uuid"
	"muserblog/models"
)

var (
	postTitle   = "Post 1"
	postContent = "Post 1 content"
)

func postForTest(user_id uuid.UUID) *models.Post {
	return &models.Post{Title: postTitle, Content: postContent, UserID: user_id}
}

func (as *ModelSuite) Test_Post() {
	prevCount := as.CountObjects(models.Post{})

	user := userForTest()
	as.DB.Create(user)

	as.NoError(as.DB.Create(postForTest(user.ID)))

	as.Equal(as.CountObjects(models.Post{})-prevCount, 1)

	user = &models.User{}
	as.DB.First(user)

	post := &models.Post{}
	as.DB.First(post)

	as.Equal(postTitle, post.Title)
	as.Equal(postContent, post.Content)
}
