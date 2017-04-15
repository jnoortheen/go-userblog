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

func countUsers(as *ModelSuite) int {
	usersCount, err := as.DB.Count(models.Users{})
	as.NoError(err)
	return usersCount
}

func (as *ModelSuite) Test_User() {
	prevCount := countUsers(as)
	user := userForTest()
	as.NoError(as.DB.Create(user))
	as.Equal(countUsers(as)-prevCount, 1)
}
