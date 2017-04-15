package models_test

import (
	"muserblog/models"
	"github.com/satori/go.uuid"
)

var (
	postTitle = "Post 1"
	postContent = "Post 1 content"
)

func postForTest(user_id uuid.UUID) *models.Post {
	return &models.Post{Title: postTitle, Content: postContent, UserID:user_id}
}

func (as *ModelSuite) Test_Post() {
	prevCount := as.countObjects(models.Post{})

	user := userForTest()
	as.DB.Create(user)

	as.NoError(as.DB.Create(postForTest(user.ID)))
	as.Equal(as.countObjects(models.Post{}) - prevCount, 1)
}
