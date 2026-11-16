package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFrequenciesRejectUnreadableSample(t *testing.T) {
	root := t.TempDir()
	for _, lang := range []string{"Python", "Ruby"} {
		dir := filepath.Join(root, lang)
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "sample"), []byte("print hello"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	file := filepath.Join(root, "Ruby", "sample")
	if err := os.Chmod(file, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(file, 0600)
	if _, err := os.ReadFile(file); err == nil {
		t.Skip("file permissions do not restrict this user")
	}
	_, err := getFrequencies(root, nil)

	if err == nil {
		t.Fatal("silently accepted an incomplete training corpus")
	}
}

func TestFrequenciesCanonicalNames(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "Fstar")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sample.fst"), []byte("module Example"), 0600); err != nil {
		t.Fatal(err)
	}
	freqs, err := getFrequencies(root, map[string]string{"Fstar": "F*"})
	if err != nil {
		t.Fatal(err)
	}
	if freqs.Languages["F*"] != 1 || len(freqs.Tokens["F*"]) == 0 || freqs.LanguageTokens["F*"] == 0 {
		t.Fatalf("missing canonical model: %+v", freqs)
	}
	if _, ok := freqs.Languages["Fstar"]; ok {
		t.Fatal("filesystem name leaked into model")
	}
}
