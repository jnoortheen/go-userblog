package models_test

import (
	"muserblog/helpers/test"
	"testing"
	"github.com/gobuffalo/suite"
)

type ModelSuite struct {
	*test.BloggySuite
}

func Test_BloggySuite(t *testing.T) {
	as := &test.BloggySuite{suite.NewModel()}
	suite.Run(t, as)
}
