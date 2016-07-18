package slinguist

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TSuite struct{}

var _ = Suite(&TSuite{})

func (s *TSuite) TestGetLanguage(c *C) {
	c.Assert(GetLanguage("foo.py", []byte{}), Equals, "Python")
	c.Assert(GetLanguage("foo.m", []byte(":- module")), Equals, "Mercury")
	c.Assert(GetLanguage("foo.m", []byte{}), Equals, "Other")
}
