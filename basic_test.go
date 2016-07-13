package slinguist

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) TestGetLanguage(c *C) {
	c.Assert(GetLanguage("foo.foo"), Equals, "Other")
	c.Assert(GetLanguage("foo.go"), Equals, "Go")
	c.Assert(GetLanguage("foo.go.php"), Equals, "PHP")
}

func (s *UtilsSuite) TestGetLanguageExtensions(c *C) {
	c.Assert(GetLanguageExtensions("foo"), HasLen, 0)
	c.Assert(GetLanguageExtensions("C"), Not(HasLen), 0)
}
