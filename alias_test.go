package slinguist

import . "gopkg.in/check.v1"

func (s *TSuite) TestGetLanguageByAlias(c *C) {
	tests := []struct {
		alias        string
		expectedLang string
	}{
		{alias: "BestLanguageEver", expectedLang: OtherLanguage},
		{alias: "aspx-vb", expectedLang: "ASP"},
		{alias: "C++", expectedLang: "C++"},
		{alias: "c++", expectedLang: "C++"},
		{alias: "objc", expectedLang: "Objective-C"},
		{alias: "golang", expectedLang: "Go"},
		{alias: "GOLANG", expectedLang: "Go"},
		{alias: "bsdmake", expectedLang: "Makefile"},
		{alias: "xhTmL", expectedLang: "HTML"},
		{alias: "python", expectedLang: "Python"},
	}

	for _, test := range tests {
		lang := GetLanguageByAlias(test.alias)
		c.Assert(lang, Equals, test.expectedLang)
	}
}
