package enry

import "testing"

func TestEnvShebangOptionOperands(t *testing.T) {
	tests := []struct{ args, want string }{
		{"-S -u PYTHONHOME python3", "Python"},
		{"-S --unset PYTHONHOME python3", "Python"},
		{"-S -C /tmp python3", "Python"},
		{"-S --chdir /tmp python3", "Python"},
		{"-S -iu PYTHONHOME python3", "Python"},
		{"-S -uPYTHONHOME python3", "Python"},
		{"-S --unset=PYTHONHOME python3", "Python"},
		{"-S -C/tmp python3", "Python"},
		{"-S -u FOO -u BAR python3 -u", "Python"},
		{"-S -u", ""},
		{"-S --unset", ""},
		{"-S -u python3", ""},
		{"-S -C python3", ""},
		{"-S -- python3 -u", "Python"},
		{"-S PYTHONHOME= python3", "Python"},
		{"-S unknown-interpreter python3", ""},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, _ := GetLanguageByShebang([]byte("#!/usr/bin/env " + tt.args))
			if got != tt.want {
				t.Fatalf("got %q; want %q", got, tt.want)
			}
		})
	}
}

func TestEnvSplitString(t *testing.T) {
	for _, args := range []string{
		`-Spython3 -u`, `--split-string=python3 -u`, `-vSpython3 -u`,
		`-S "python3" -u`, `-S 'python3' -u`, `-S -u "MY VARIABLE" python3`,
		`-S python3\_-u`, `-S python3 # comment`, `-S python3\c ignored`,
	} {
		t.Run(args, func(t *testing.T) {
			got, _ := GetLanguageByShebang([]byte("#!/usr/bin/env " + args))
			if got != "Python" {
				t.Fatalf("got %q", got)
			}
		})
	}
	for _, args := range []string{`-S "python3`, `-S 'python3`, `-S ${UNKNOWN}`, `-S python3\q`, `-S "python3 -u"`} {
		got, _ := GetLanguageByShebang([]byte("#!/usr/bin/env " + args))
		if got != "" {
			t.Errorf("%q: got %q", args, got)
		}
	}
}
