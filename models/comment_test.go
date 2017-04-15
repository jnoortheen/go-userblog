package models_test

import (
	"muserblog/models"
)

var (
	commentContent = "test_user"
)

func (as *ModelSuite) Test_Comment() {
	prevCount := as.countObjects(models.Comment{})

	user := userForTest()
	as.NoError(as.DB.Create(user))

	post := postForTest()
	as.NoError(as.DB.Create(post))

	comment := &models.Comment{Content:commentContent, PostID:post.ID, UserID:user.ID}
	as.NoError(as.DB.Create(comment))

	as.Equal(as.countObjects(models.Comment{}) - prevCount, 1)
}
