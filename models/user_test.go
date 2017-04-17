package models_test

import (
	"muserblog/models"
	//"fmt"
)

var (
	userName = "test_user"
	userPwdPlain = "password"
	userEmail = "user@mail.com"
)

func userForTest() *models.User {
	return models.NewUser(userName, userPwdPlain, userEmail)
}

func (as *ModelSuite) Test_User() {
	prevCount := as.CountObjects(models.User{})

	user := userForTest()
	as.NoValidationError(as.DB.ValidateAndCreate(user))

	as.Equal(as.CountObjects(models.User{}) - prevCount, 1)

	user = userForTest()
	verr, err := as.DB.ValidateAndCreate(user)
	as.HasValidationError(verr, err, models.UniqUserNameErrMsg)
}
