package enry

import (
	"reflect"
	"testing"
)

func TestClassifierTiesAreDeterministic(t *testing.T) {
	c := &naiveBayes{
		languagesLogProbabilities: map[string]float64{"Python": -1, "Ruby": -1, "Go": -0.5},
	}
	for i := 0; i < 100; i++ {
		got := c.classify(nil, map[string]float64{"Ruby": 1, "Python": 1, "Go": 1})
		want := []string{"Go", "Python", "Ruby"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}
	}
}

func TestClassifierFilesystemAlias(t *testing.T) {
	for _, name := range []string{"F*", "Fstar"} {
		got, safe := GetLanguageByClassifier([]byte("module Example\nlet x = 1"), []string{name})
		if got != "F*" || !safe {
			t.Fatalf("candidate %q: got %q, safe=%v", name, got, safe)
		}
	}
}
