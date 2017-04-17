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

	commentUser := userForTest()
	as.NoError(as.DB.Create(commentUser))

	post := postForTest(user.ID)
	as.NoError(as.DB.Create(post))

	comment := &models.Comment{Content: commentContent, PostID: post.ID, UserID: commentUser.ID}
	as.NoError(as.DB.Create(comment))

	as.Equal(as.countObjects(models.Comment{}) - prevCount, 1)

	comment = &models.Comment{}
	as.DB.First(comment)
	as.Equal(commentContent, comment.Content)
}
