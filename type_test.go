package slinguist

import . "gopkg.in/check.v1"

func (s *TSuite) TestGetLanguageType(c *C) {
	langType := GetLanguageType("BestLanguageEver")
	c.Assert(langType, Equals, Unknown)

	langType = GetLanguageType("JSON")
	c.Assert(langType, Equals, Data)

	langType = GetLanguageType("COLLADA")
	c.Assert(langType, Equals, Data)

	langType = GetLanguageType("Go")
	c.Assert(langType, Equals, Programming)

	langType = GetLanguageType("Brainfuck")
	c.Assert(langType, Equals, Programming)

	langType = GetLanguageType("HTML")
	c.Assert(langType, Equals, Markup)

	langType = GetLanguageType("Sass")
	c.Assert(langType, Equals, Markup)

	langType = GetLanguageType("AsciiDoc")
	c.Assert(langType, Equals, Prose)

	langType = GetLanguageType("Textile")
	c.Assert(langType, Equals, Prose)
}
