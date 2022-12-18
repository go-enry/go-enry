// Copyright 2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://cs.opensource.google/go/x/perf/+/e8d778a6:LICENSE

package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// text_diff returns a human-readable description of the differences between s1 and s2.
// It essentially is https://cs.opensource.google/go/x/perf/+/e8d778a6:internal/diff/diff.go
// Used only in code generator tests, as a debugging aid.
// It is not part of any release artifact and is not distibuted with enry.
func text_diff(b1, b2 []byte) (string, error) {
	if bytes.Equal(b1, b2) {
		return "", nil
	}

	cmd := "diff"
	if _, err := exec.LookPath(cmd); err != nil {
		return "", fmt.Errorf("diff command unavailable\nold: %q\nnew: %q", b1, b2)
	}

	f1, err := writeTempFile("", "gen_test", b1)
	if err != nil {
		return "", err
	}
	defer os.Remove(f1)

	f2, err := writeTempFile("", "gen_test", b2)
	if err != nil {
		return "", err
	}
	defer os.Remove(f2)

	data, err := exec.Command(cmd, "-u", f1, f2).CombinedOutput()
	if len(data) > 0 { // diff exits with a non-zero status when the files don't match
		err = nil
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeTempFile(dir, prefix string, data []byte) (string, error) {
	file, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}
