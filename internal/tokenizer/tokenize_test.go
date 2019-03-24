package tokenizer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testContent = `#!/usr/bin/ruby

#!/usr/bin/env node

aaa

#!/usr/bin/env A=B foo=bar awk -f

#!python

func Tokenize(content []byte) []string {
	splitted := bytes.Fields(content)
	tokens := /* make([]string, 0, len(splitted))
	no comment -- comment
	for _, tokenByte := range splitted {
		token64 := base64.StdEncoding.EncodeToString(tokenByte)
		tokens = append(tokens, token64)
		notcatchasanumber3.5
	}*/
othercode
	/* testing multiple 
	
		multiline comments*/

<!-- com
	ment -->
<!-- comment 2-->
ppp no comment # comment

"literal1"

abb (tokenByte, 0xAF02) | ,3.2L

'literal2' notcatchasanumber3.5

	5 += number * anotherNumber
	if isTrue && isToo {
		0b00001000 >> 1
	}

	return tokens

oneBool = 3 <= 2
varBool = 3<=2>
 
#ifndef
#i'm not a comment if the single line comment symbol is not followed by a white

  PyErr_SetString(PyExc_RuntimeError, "Relative import is not supported for Python <=2.4.");

<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
    <head>
        <title id="hola" class="">This is a XHTML sample file</title>
        <style type="text/css"><![CDATA[
            #example {
                background-color: yellow;
            }
        ]]></style>
    </head>
    <body>
        <div id="example">
            Just a simple <strong>XHTML</strong> test page.
        </div>
    </body>
</html>`
)

var (
	tokensFromTestContent = []string{"SHEBANG#!ruby", "SHEBANG#!node", "SHEBANG#!awk", "<!DOCTYPE>", "PUBLIC", "W3C", "DTD", "XHTML", "1", "0",
		"Strict", "EN", "http", "www", "w3", "org", "TR", "xhtml1", "DTD", "xhtml1", "strict", "dtd", "<html>", "<head>", "<title>", "class=",
		"</title>", "<style>", "<![CDATA[>", "example", "background", "color", "yellow", "</style>", "</head>", "<body>", "<div>", "<strong>",
		"</strong>", "</div>", "</body>", "</html>", "(", "[", "]", ")", "[", "]", "{", "(", ")", "(", ")", "{", "}", "(", ")", ";", "{", ";",
		"}", "]", "]", "#", "/usr/bin/ruby", "#", "/usr/bin/env", "node", "aaa", "#", "/usr/bin/env", "A", "B", "foo", "bar", "awk", "f", "#",
		"python", "func", "Tokenize", "content", "byte", "string", "splitted", "bytes.Fields", "content", "tokens", "othercode", "ppp", "no",
		"comment", "abb", "tokenByte", "notcatchasanumber", "number", "*", "anotherNumber", "if", "isTrue", "isToo", "b", "return", "tokens",
		"oneBool", "varBool", "#ifndef", "#i", "m", "not", "a", "comment", "if", "the", "single", "line", "comment", "symbol", "is", "not",
		"followed", "by", "a", "white", "PyErr_SetString", "PyExc_RuntimeError", "html", "PUBLIC", "xmlns", "id", "class", "This", "is", "a",
		"XHTML", "sample", "file", "type", "#example", "background", "color", "yellow", "id", "Just", "a", "simple", "XHTML", "test", "page.",
		"-", "|", "+", "&&", "<", "<", "-", "!", "!", "!", "=", "=", "!", ":", "=", ":", "=", ",", ",", "=", ">", ">", "=", "=", "=", "=", ">",
		"'", ",", ">", "=", ">", "=", "=", ">", "=", ">", ":", ">", "=", ">"}

	tests = []struct {
		name     string
		content  []byte
		expected []string
	}{
		{name: "content", content: []byte(testContent), expected: tokensFromTestContent},
	}
)

func TestTokenize(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := string(test.content)
			tokens := TokenizeFlex(test.content)
			after := string(test.content)
			require.Equal(t, before, after, "the input slice was modified")
			require.Equal(t, len(test.expected), len(tokens), fmt.Sprintf("token' slice length = %v, want %v", len(test.expected), len(tokens)))
			for i, expectedToken := range test.expected {
				assert.Equal(t, expectedToken, tokens[i], fmt.Sprintf("token = %v, want %v", tokens[i], expectedToken))
			}
		})
	}
}

func BenchmarkTokenizer_BaselineCopy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			test.content = append([]byte(nil), test.content...)
		}
	}
}

func BenchmarkTokenizerGo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			Tokenize(test.content)
		}
	}
}

func BenchmarkTokenizerC(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			TokenizeC(test.content)
		}
	}
}

func BenchmarkTokenizerFlex(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			TokenizeFlex(test.content)
		}
	}
}

//TODO(bzz): introduce tokenizer benchmark suit
// baseline - just read the files
// RE2
// oniguruma
// cgo to flex-based impl
