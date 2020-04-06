package data

import "gopkg.in/toqueteos/substring.v1"

// TestMatchers is hand made collection of regexp used by the function `enry.IsTest`
// to identify test files in different languages.
var TestMatchers = substring.Or(
	substring.Regexp(`(^|/)tests/.*Test\.php$`),
	substring.Regexp(`(^|/)test/.*Test(s?)\.java$`),
	substring.Regexp(`(^|/)test(/|/.*/)Test.*\.java$`),
	substring.Regexp(`(^|/)test/.*(Test(s?)|Spec(s?))\.scala$`),
	substring.Regexp(`(^|/)test_.*\.py$`),
	substring.Regexp(`(^|/).*_test\.go$`),
	substring.Regexp(`(^|/).*_(test|spec)\.rb$`),
	substring.Regexp(`(^|/).*Test(s?)\.cs$`),
	substring.Regexp(`(^|/).*\.(test|spec)\.(ts|tsx|js)$`),
)
