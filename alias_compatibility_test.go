package enry_test

import (
	"testing"

	enry "github.com/go-enry/go-enry/v2"
)

func TestHyphenatedLanguageAliases(t *testing.T) {
	for _, test := range []struct {
		alias string
		want  string
	}{
		{"OpenStep-Property-List", "OpenStep Property List"},
		{"openstep_property_list", "OpenStep Property List"},
		{"OpenStep Property List", "OpenStep Property List"},
		{"POV-Ray-SDL", "POV-Ray SDL"},
		{"pov-ray_sdl", "POV-Ray SDL"},
		{"Tree-sitter-Query", "Tree-sitter Query"},
		{"tree-sitter_query", "Tree-sitter Query"},
		{"coffee-script", "CoffeeScript"},
		{"aspx-vb", "ASP.NET"},
		{"OpenStep-Property-List,ignored", "OpenStep Property List"},
	} {
		t.Run(test.alias, func(t *testing.T) {
			got, ok := enry.GetLanguageByAlias(test.alias)
			if !ok || got != test.want {
				t.Fatalf("GetLanguageByAlias(%q) = %q, %v; want %q, true", test.alias, got, ok, test.want)
			}
		})
	}
}
