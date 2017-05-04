package slinguist

import . "gopkg.in/check.v1"

func (s *TSuite) TestGetLanguageByFilename(c *C) {
	lang, safe := GetLanguageByFilename(`unknown.interpreter`)
	c.Assert(lang, Equals, OtherLanguage)
	c.Assert(safe, Equals, false)

	lang, safe = GetLanguageByFilename(`.bashrc`)
	c.Assert(lang, Equals, "Shell")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByFilename(`Dockerfile`)
	c.Assert(lang, Equals, "Dockerfile")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByFilename(`Makefile.frag`)
	c.Assert(lang, Equals, "Makefile")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByFilename(`makefile`)
	c.Assert(lang, Equals, "Makefile")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByFilename(`Vagrantfile`)
	c.Assert(lang, Equals, "Ruby")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByFilename(`_vimrc`)
	c.Assert(lang, Equals, "Vim script")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByFilename(`pom.xml`)
	c.Assert(lang, Equals, "Maven POM")
	c.Assert(safe, Equals, true)
}
