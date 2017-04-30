package models_test

import (
	"github.com/gobuffalo/suite"
	"testing"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	as := &ModelSuite{suite.NewModel()}
	suite.Run(t, as)
}

// count and assert no error happens
func (as *ModelSuite) CountObjects(model interface{}) int {
	cnt, err := as.DB.Count(model)
	as.NoError(err)
	return cnt
}
