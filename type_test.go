package slinguist

import . "gopkg.in/check.v1"

func (s *TSuite) TestGetLanguageType(c *C) {
	langType := GetLanguageType("BestLanguageEver")
	c.Assert(langType, Equals, TypeUnknown)

	langType = GetLanguageType("JSON")
	c.Assert(langType, Equals, TypeData)

	langType = GetLanguageType("COLLADA")
	c.Assert(langType, Equals, TypeData)

	langType = GetLanguageType("Go")
	c.Assert(langType, Equals, TypeProgramming)

	langType = GetLanguageType("Brainfuck")
	c.Assert(langType, Equals, TypeProgramming)

	langType = GetLanguageType("HTML")
	c.Assert(langType, Equals, TypeMarkup)

	langType = GetLanguageType("Sass")
	c.Assert(langType, Equals, TypeMarkup)

	langType = GetLanguageType("AsciiDoc")
	c.Assert(langType, Equals, TypeProse)

	langType = GetLanguageType("Textile")
	c.Assert(langType, Equals, TypeProse)
}
