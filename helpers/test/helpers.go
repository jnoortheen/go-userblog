package test

import (
	"github.com/gobuffalo/suite"
	"github.com/markbates/validate"
)

type BloggySuite struct {
	*suite.Model
}


// count and assert no error happens
func (as *BloggySuite) CountObjects(model interface{}) int {
	usersCount, err := as.DB.Count(model)
	as.NoError(err)
	return usersCount
}

// assert no validation error happended
func (bs *BloggySuite) NoValidationError(verrs *validate.Errors, err error) {
	bs.NoError(err)
	bs.NoError(verrs)
}

// assert no validation error happended
func (bs *BloggySuite) HasValidationError(verrs *validate.Errors, err error, expectedValidationErrMsg string) {
	bs.Error(err)
	bs.EqualError(verrs, expectedValidationErrMsg)
}
