package generator

import "testing"

func TestBuildAliasLanguageMapHyphens(t *testing.T) {
	languages := map[string]*languageInfo{
		"Example Language": {Aliases: []string{"other name", "explicit-alias"}},
		"POV-Ray SDL":      {},
	}
	got := buildAliasLanguageMap(languages, getAlphabeticalOrderedKeys(languages))
	for alias, want := range map[string]string{
		"example_language": "Example Language",
		"example-language": "Example Language",
		"other_name":       "Example Language",
		"explicit-alias":   "Example Language",
		"pov-ray_sdl":      "POV-Ray SDL",
		"pov-ray-sdl":      "POV-Ray SDL",
	} {
		if got[alias] != want {
			t.Errorf("alias %q = %q; want %q", alias, got[alias], want)
		}
	}
	if _, exists := got["other-name"]; exists {
		t.Error("explicit aliases must not gain additional hyphenated spellings")
	}
}

func TestBuildAliasLanguageMapPreservesExistingAliases(t *testing.T) {
	languages := map[string]*languageInfo{
		"First":      {Aliases: []string{"other-name"}},
		"Other Name": {},
	}
	got := buildAliasLanguageMap(languages, getAlphabeticalOrderedKeys(languages))
	if got["other-name"] != "First" || got["other_name"] != "Other Name" {
		t.Fatalf("new hyphenated names must not replace existing aliases: %v", got)
	}
}
