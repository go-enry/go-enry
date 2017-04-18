package slinguist

import . "gopkg.in/check.v1"

const (
	multilineExecHack = `#!/bin/sh
# Next line is comment in Tcl, but not in sh... \
exec tclsh "$0" ${1+"$@"}`

	multilineNoExecHack = `#!/bin/sh
#<<<#
echo "A shell script in a zkl program ($0)"
echo "Now run zkl <this file> with Hello World as args"
zkl $0 Hello World!
exit
#<<<#
println("The shell script says ",vm.arglist.concat(" "));`
)

func (s *TSuite) TestGetLanguageByShebang(c *C) {
	lang, safe := GetLanguageByShebang([]byte(`#!/unknown/interpreter`))
	c.Assert(lang, Equals, OtherLanguage)
	c.Assert(safe, Equals, false)

	lang, safe = GetLanguageByShebang([]byte(`no shebang`))
	c.Assert(lang, Equals, OtherLanguage)
	c.Assert(safe, Equals, false)

	lang, safe = GetLanguageByShebang([]byte(`#!/usr/bin/env`))
	c.Assert(lang, Equals, OtherLanguage)
	c.Assert(safe, Equals, false)

	lang, safe = GetLanguageByShebang([]byte(`#!/usr/bin/python -tt`))
	c.Assert(lang, Equals, "Python")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByShebang([]byte(`#!/usr/bin/env python2.6`))
	c.Assert(lang, Equals, "Python")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByShebang([]byte(`#!/usr/bin/env perl`))
	c.Assert(lang, Equals, "Perl")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByShebang([]byte(`#!	/bin/sh`))
	c.Assert(lang, Equals, "Shell")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByShebang([]byte(`#!bash`))
	c.Assert(lang, Equals, "Shell")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByShebang([]byte(multilineExecHack))
	c.Assert(lang, Equals, "Tcl")
	c.Assert(safe, Equals, true)

	lang, safe = GetLanguageByShebang([]byte(multilineNoExecHack))
	c.Assert(lang, Equals, "Shell")
	c.Assert(safe, Equals, true)
}
