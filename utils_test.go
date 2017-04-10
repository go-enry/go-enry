package slinguist

import (
	"bytes"

	. "gopkg.in/check.v1"
)

func (s *TSuite) TestIsVendor(c *C) {
	c.Assert(IsVendor("foo/bar"), Equals, false)
	c.Assert(IsVendor("foo/vendor/foo"), Equals, true)
	c.Assert(IsVendor(".sublime-project"), Equals, true)
	c.Assert(IsVendor("leaflet.draw-src.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/MochiKit.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/dojo.js"), Equals, true)
	c.Assert(IsVendor("foo/env/whatever"), Equals, true)
	c.Assert(IsVendor("foo/.imageset/bar"), Equals, true)
	c.Assert(IsVendor("Vagrantfile"), Equals, true)
}

func (s *TSuite) TestIsDocumentation(c *C) {
	c.Assert(IsDocumentation("foo"), Equals, false)
	c.Assert(IsDocumentation("README"), Equals, true)
}

func (s *TSuite) TestIsConfiguration(c *C) {
	c.Assert(IsConfiguration("foo"), Equals, false)
	c.Assert(IsConfiguration("foo.ini"), Equals, true)
	c.Assert(IsConfiguration("foo.json"), Equals, true)
}

func (s *TSuite) TestIsBinary(c *C) {
	c.Assert(IsBinary([]byte("foo")), Equals, false)

	binary := []byte{0}
	c.Assert(IsBinary(binary), Equals, true)

	binary = bytes.Repeat([]byte{'o'}, 8000)
	binary = append(binary, byte(0))
	c.Assert(IsBinary(binary), Equals, false)
}

const (
	htmlPath = "some/random/dir/file.html"
	jsPath   = "some/random/dir/file.js"
)

func (s *TSuite) BenchmarkVendor(c *C) {
	for i := 0; i < c.N; i++ {
		_ = IsVendor(htmlPath)
	}
}

func (s *TSuite) BenchmarkVendorJS(c *C) {
	for i := 0; i < c.N; i++ {
		_ = IsVendor(jsPath)
	}
}
