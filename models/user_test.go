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
	prevCount := as.CountObjects(models.User{})

	user := userForTest()
	verr, err := as.DB.ValidateAndCreate(user)
	as.NoError(err)
	if verr.HasAny() {
		as.Fail(verr.Error())
	}
	as.Equal(as.CountObjects(models.User{})-prevCount, 1)

	user = userForTest()
	verr, err = as.DB.ValidateAndCreate(user)
	as.EqualError(verr, models.UniqUserNameErrMsg)
}

func (as *ModelSuite) Test_UserHashFunction() {
	user := userForTest()
	verr, err := as.DB.ValidateAndCreate(user)
	as.NoError(err)
	if verr.HasAny() {
		as.Fail(verr.Error())
	}
	as.Equal(user.AuthToken(), user.AuthToken())
}
