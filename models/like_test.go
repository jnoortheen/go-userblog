package models_test

import "muserblog/models"

func (as *ModelSuite) Test_Like() {
	prevCount := as.countObjects(models.Like{})

	user := userForTest()
	as.NoError(as.DB.Create(user))

	post := postForTest(user.ID)
	as.NoError(as.DB.Create(post))

	comment := &models.Like{PostID: post.ID, UserID: user.ID}
	as.NoError(as.DB.Create(comment))

	as.Equal(as.countObjects(models.Like{}) - prevCount, 1)
}
