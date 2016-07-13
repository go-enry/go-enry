package slinguist

import . "gopkg.in/check.v1"

func (s *TSuite) TestGetLanguageByExtension(c *C) {
	lang, safe := GetLanguageByExtension("foo.foo")
	c.Assert(lang, Equals, "Other")
	c.Assert(safe, Equals, false)

	lang, safe = GetLanguageByExtension("foo.go")
	c.Assert(lang, Equals, "Go")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByExtension("foo.go.php")
	c.Assert(lang, Equals, "PHP")
	c.Assert(safe, Equals, false)
}

func (s *TSuite) TestGetLanguageExtensions(c *C) {
	c.Assert(GetLanguageExtensions("foo"), HasLen, 0)
	c.Assert(GetLanguageExtensions("C"), Not(HasLen), 0)
}
