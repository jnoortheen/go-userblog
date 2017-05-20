package models_test

import (
	"muserblog/models"
)

var (
	commentContent = "content"
)

func (as *ModelSuite) Test_Comment() {
	prevCount := as.CountObjects(models.Comment{})

	user := userForTest()
	as.NoError(as.DB.Create(user))

	commentUser := userForTest()
	as.NoError(as.DB.Create(commentUser))

	post := postForTest(user.ID)
	as.NoError(as.DB.Create(post))

	comment := &models.Comment{Content: commentContent, PostID: post.ID, UserID: commentUser.ID}
	as.NoError(as.DB.Create(comment))

	as.Equal(as.CountObjects(models.Comment{})-prevCount, 1)

	comment = &models.Comment{}
	as.DB.First(comment)
	as.Equal(commentContent, comment.Content)
}

func (as *ModelSuite) Test_CommentExt() {
	user := userForTest()
	as.NoError(as.DB.Create(user))

	post := postForTest(user.ID)
	as.NoError(as.DB.Create(post))

	comment := &models.Comment{Content: commentContent, PostID: post.ID, UserID: user.ID}
	as.NoError(as.DB.Create(comment))
	commentExt := &models.CommentExt{}
	commentExt.Comment = comment
	commentExt.Update(as.DB, user)
	for _, cont := range []string{user.Name, user.ID.String(), post.ID.String(), commentContent} {
		as.Contains(commentExt.String(), cont)
	}
}
