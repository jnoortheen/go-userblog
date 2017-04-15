package models_test

import (
	"testing"

	"github.com/gobuffalo/suite"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	as := &ModelSuite{suite.NewModel()}
	suite.Run(t, as)
}

// count and assert no error happens
func (as *ModelSuite) countObjects(model interface{}) int {
	usersCount, err := as.DB.Count(model)
	as.NoError(err)
	return usersCount
}
