package slinguist

import . "gopkg.in/check.v1"

func (s *TSuite) TestGetLanguageByAlias(c *C) {
	tests := []struct {
		alias        string
		expectedLang string
		expectedOk   bool
	}{
		{alias: "BestLanguageEver", expectedLang: OtherLanguage, expectedOk: false},
		{alias: "aspx-vb", expectedLang: "ASP", expectedOk: true},
		{alias: "C++", expectedLang: "C++", expectedOk: true},
		{alias: "c++", expectedLang: "C++", expectedOk: true},
		{alias: "objc", expectedLang: "Objective-C", expectedOk: true},
		{alias: "golang", expectedLang: "Go", expectedOk: true},
		{alias: "GOLANG", expectedLang: "Go", expectedOk: true},
		{alias: "bsdmake", expectedLang: "Makefile", expectedOk: true},
		{alias: "xhTmL", expectedLang: "HTML", expectedOk: true},
		{alias: "python", expectedLang: "Python", expectedOk: true},
	}

	for _, test := range tests {
		lang, ok := GetLanguageByAlias(test.alias)
		c.Assert(lang, Equals, test.expectedLang)
		c.Assert(ok, Equals, test.expectedOk)
	}
}
