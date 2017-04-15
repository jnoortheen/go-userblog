package models_test

import (
	"muserblog/models"
)

var (
	userName     = "test_user"
	userPwdPlain = "password"
	userEmail    = "user@mail.com"
)

func userForTest() *models.User {
	return models.NewUser(userName, userPwdPlain, userEmail)
}

func (as *ModelSuite) Test_User() {
	prevCount := as.countObjects(models.User{})
	user := userForTest()
	as.NoError(as.DB.Create(user))
	as.Equal(as.countObjects(models.User{})-prevCount, 1)
}
