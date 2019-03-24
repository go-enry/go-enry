// +build flex

package tokenizer

import "gopkg.in/src-d/enry.v1/internal/tokenizer/flex"

// Tokenize returns language-agnostic lexical tokens from content. The tokens
// returned should match what the Linguist library returns. At most the first
// 100KB of content are tokenized.
func Tokenize(content []byte) []string {
	if len(content) > byteLimit {
		content = content[:byteLimit]
	}

	return flex.TokenizeFlex(content)
}
