package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Re-used by the packages: enry (test), enry (benchmark) and code-generator (test).
// Does not rely on testify, panics on errors so that there always is a trace to the caller.
func MaybeCloneLinguist(envVar, url, commit string) (string, bool, error) {
	var err error
	linguistTmpDir := os.Getenv(envVar)
	isCleanupNeeded := false
	isLinguistCloned := linguistTmpDir != ""
	if !isLinguistCloned {
		linguistTmpDir, err = ioutil.TempDir("", "linguist-")
		if err != nil {
			return "", false, err
		}

		isCleanupNeeded = true
		cmd := exec.Command("git", "clone", "--depth", "100", url, linguistTmpDir)
		if err := cmd.Run(); err != nil {
			panic(fmt.Errorf("%s: %w", cmd.String(), err))
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err = os.Chdir(linguistTmpDir); err != nil {
		panic(err)
	}

	cmd := exec.Command("git", "checkout", commit)
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("%s: %w", cmd.String(), err))
	}

	if err = os.Chdir(cwd); err != nil {
		panic(fmt.Errorf("%s: %w", cmd.String(), err))
	}
	return linguistTmpDir, isCleanupNeeded, nil
}
