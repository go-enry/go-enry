# `.gitattributes` Edge-Case Testing Report

**PR Under Review:** [#206 — Add .gitattributes support for language detection overrides](https://github.com/go-enry/go-enry/pull/206)
**Branch tested:** `claude/test-gitattributes-edge-cases-T0VNd` (rebased onto `refs/pull/206/head` @ `90d6048`)
**Test file:** `gitattributes_edge_test.go`
**Date:** 2026-03-20

---

## Executive Summary

PR #206 adds Linguist-compatible `.gitattributes` support with a well-structured API. All original PR tests pass. After the first round of review, two bugs were found and fixed in commits `94aca8b` (glob fixes) and `90d6048` (`IsDetectable` tri-state). After rebasing onto those fixes and re-running the full edge-case suite (106+ library tests + 10 CLI tests), **one remaining bug** was discovered in the trailing-slash fix.

| Category | Round 1 | Round 2 (latest) |
|---|---|---|
| Confirmed bugs (wrong output) | 2 | **1 remaining** |
| Bugs fixed since round 1 | — | 2 ✅ |
| Behavioral quirks documented | 4 | 4 |
| Total edge cases | 106+ lib + 10 CLI | 108+ lib + 10 CLI |

---

## Bugs Fixed Since Round 1

### ✅ Fixed: `**` did not enforce path-component boundaries (commit `94aca8b`)

**Was:** `**/vendor/**` matched `notvendor/lib.go` because `globMatch` tried every byte offset.

**Fix:** The `**` loop now only tries at `/`-aligned boundaries:
```go
for i := 0; i <= len(str); i++ {
    if i == 0 || str[i-1] == '/' {
        if globMatch(rest, str[i:]) { return true }
    }
}
```
**CLI verification:** `notvendor/lib.go` now correctly appears in output with `**/vendor/** linguist-vendored`. ✅

---

### ✅ Fixed: Trailing-slash pattern did not match direct file children (commit `94aca8b`)

**Was:** `vendor/foo.go` was not matched by `vendor/`.

**Fix:** After stripping the trailing slash from the pattern, files whose path starts with `dir+"/`" are matched immediately:
```go
} else if strings.HasPrefix(path, dir+"/") {
    return true
}
```
**Test status:** `vendor/foo.go` and `vendor/lib/bar.go` now match `vendor/`. ✅

---

### ✅ Fixed: `IsDetectable` now returns tri-state `(bool, bool)` (commit `90d6048`)

`IsDetectable(path string) (bool, bool)` now returns `(value, hasOverride)`, matching `GetLanguage`'s pattern. The CLI uses this to properly:
- Include data/prose languages forced via `linguist-detectable`
- Exclude programming languages forced out via `-linguist-detectable`

---

## Remaining Bug

### Bug 3: Trailing-slash pattern does not match nested subdirectories

**Severity:** HIGH — causes `vendor/` patterns to silently fail for nested dependency trees

**Location:** `gitattributes.go` → `matchGitPattern()`, trailing-slash branch

**Root Cause:**
The trailing-slash fix handles three cases:
1. Path ending with `/` (directory): strips trailing slash, then does basename glob matching
2. Path not ending with `/` that starts with `dir+"/`": returns `true`
3. Everything else: returns `false`

Case 1 is the problem. For path `vendor/github.com/` with pattern `vendor/`:
- `dir = "vendor"`, path ends with `/`
- Code strips trailing slash → `path = "vendor/github.com"`, `pattern = "vendor"`
- No slash in pattern → basename match: `filepath.Base("vendor/github.com") = "github.com"`
- `globMatch("vendor", "github.com")` → **false**
- Falls back to default `IsVendor` → **true** → directory is `SkipDir`'d

**Real-world impact (confirmed):**
The CLI's `filepath.Walk` adds a trailing `/` to directory paths. When it walks `vendor/`, it encounters intermediate directories like `vendor/github.com/` and `vendor/github.com/pkg/` as separate walk entries with trailing slashes. The pattern `vendor/` only un-vendors the immediate directory path itself and its **direct file children**, not its subdirectory tree.

```
# .gitattributes
vendor/ -linguist-vendored

# Expected: vendor/github.com/pkg/errors.go counted in stats
# Actual:   vendor/github.com/ is still IsVendor=true → SkipDir
#           → errors.go is never reached by the walker
```

**Failing test cases (new, confirmed):**

| Pattern | Path | Expected | Got |
|---|---|---|---|
| `vendor/` | `vendor/github.com/` | `true` | `false` ❌ |
| `vendor/` | `vendor/github.com/pkg/` | `true` | `false` ❌ |

**Workaround:** Use `vendor/**` instead of `vendor/` in `.gitattributes`. This already works correctly.

**Fix:** In the trailing-slash branch, check both direct and nested subdirectory paths:
```go
if strings.HasSuffix(path, "/") {
    // For directory paths: match if it is the directory itself OR a subdirectory
    dirPath := strings.TrimSuffix(path, "/")
    if dirPath == dir || strings.HasPrefix(dirPath, dir+"/") {
        return true
    }
    return false
}
```

---

## Behavioral Quirks (unchanged from Round 1)

### Quirk 1: Inline comments are not supported
`#` mid-line is parsed as an attribute name, not a comment. `*.go linguist-vendored # reason` creates attrs for `#`, `reason`, etc. Correct per git spec but may surprise users.

### Quirk 2: `!attr` resets attribute (new in `94aca8b`)
The PR now handles `!linguist-vendored` as "unspecified" (reset to default). `getAttr` returns `("", false)` for unspecified attributes, causing the method to fall back to default detection. This is correct Linguist behavior.

### Quirk 3: Unknown language names returned verbatim
`GetLanguage` returns raw unrecognized language names (e.g. `"GoLangg"`) with `ok=true`. No warning is logged.

### Quirk 4: Escaped spaces in filenames silently break patterns
`strings.Fields` splits on unescaped spaces. A filename like `path\ with\ spaces` causes incorrect parsing. Documented but not yet addressed.

---

## Comprehensive Edge-Case Results (Round 2)

### Pattern Matching (`matchGitPattern`)

| # | Pattern | Path | Expected | Result |
|---|---|---|---|---|
| 1 | `**/foo.go` | `barfoo.go` | `false` | ✅ PASS (fixed) |
| 2 | `**/foo.go` | `xfoo.go` | `false` | ✅ PASS (fixed) |
| 3 | `**/test` | `atest` | `false` | ✅ PASS (fixed) |
| 4 | `**/vendor/**` | `bvendor/lib.go` | `false` | ✅ PASS (fixed) |
| 5 | `**/vendor/**` | `notvendor/x.go` | `false` | ✅ PASS (fixed) |
| 6 | `**/foo.go` | `foo.go` | `true` | ✅ PASS |
| 7 | `**/foo.go` | `src/foo.go` | `true` | ✅ PASS |
| 8 | `**/foo.go` | `src/lib/foo.go` | `true` | ✅ PASS |
| 9 | `**/vendor/**` | `vendor/foo.go` | `true` | ✅ PASS |
| 10 | `**/vendor/**` | `a/vendor/foo.go` | `true` | ✅ PASS |
| 11 | `vendor/` | `vendor/foo.go` | `true` | ✅ PASS (fixed) |
| 12 | `vendor/` | `vendor/lib/bar.go` | `true` | ✅ PASS (fixed) |
| 13 | `vendor/` | `vendor/` | `true` | ✅ PASS |
| 14 | `vendor/` | `vendor/github.com/` | `true` | ❌ FAIL (Bug #3) |
| 15 | `vendor/` | `vendor/github.com/pkg/` | `true` | ❌ FAIL (Bug #3) |
| 16 | `vendor/` | `vendor` | `false` | ✅ PASS |
| 17 | `vendor/` | `src/vendor/foo.go` | `false` | ✅ PASS |
| 18 | `vendor/` | `notvendor/foo.go` | `false` | ✅ PASS |
| 19 | `**` | `a/b/c.go` | `true` | ✅ PASS |
| 20 | `src/**/test.go` | `src/test.go` | `true` | ✅ PASS |
| 21 | `src/**/test.go` | `src/a/b/test.go` | `true` | ✅ PASS |
| 22 | `src/**/test.go` | `other/a/test.go` | `false` | ✅ PASS |
| 23 | `a/**/b/**/c.go` | `a/b/c.go` | `true` | ✅ PASS |
| 24 | `a/**/b/**/c.go` | `a/x/b/y/c.go` | `true` | ✅ PASS |
| 25 | `a/**/b/**/c.go` | `a/b/x/c.go` | `true` | ✅ PASS |
| 26 | `a/**/b/**/c.go` | `a/x/c.go` | `false` | ✅ PASS |
| 27 | `*/*.go` | `src/main.go` | `true` | ✅ PASS |
| 28 | `*/*.go` | `a/b/main.go` | `false` | ✅ PASS |
| 29 | `src/*.go` | `src/sub/main.go` | `false` | ✅ PASS |
| 30 | `*.go` | `a/b/c.go` | `true` | ✅ PASS (basename) |
| 31 | `?.go` | `a.go` | `true` | ✅ PASS |
| 32 | `?.go` | `ab.go` | `false` | ✅ PASS |
| 33 | `?.go` | `/.go` | `false` | ✅ PASS |
| 34 | `/Makefile` | `Makefile` | `true` | ✅ PASS |
| 35 | `/Makefile` | `src/Makefile` | `false` | ✅ PASS |
| 36 | `Makefile` | `a/b/c/Makefile` | `true` | ✅ PASS |
| 37 | `/vendor/**` | `vendor/a/b.go` | `true` | ✅ PASS |
| 38 | `/vendor/**` | `src/vendor/a.go` | `false` | ✅ PASS |
| 39 | `*.[ch]` | `foo.c` | `true` | ✅ PASS |
| 40 | `*.[ch]` | `foo.go` | `false` | ✅ PASS |
| 41 | `*.[^ch]` | `foo.c` | `false` | ✅ PASS |
| 42 | `*.[!ch]` | `foo.c` | `false` | ✅ PASS |
| 43 | `*.[a-z]` | `foo.c` | `true` | ✅ PASS |
| 44 | `*.[A-Z]` | `foo.C` | `true` | ✅ PASS |
| 45 | `[a-z]oo.go` | `/oo.go` | `false` | ✅ PASS |
| 46 | `*.c++` | `main.c++` | `true` | ✅ PASS (literal) |
| 47 | `*.c++` | `main.cpp` | `false` | ✅ PASS |
| 48 | `foo.go` | `fooXgo` | `false` | ✅ PASS (dot is literal) |
| 49 | `**/*.go` | `a/b/c/d/e/f/g/h/i/j/main.go` | `true` | ✅ PASS |
| 50–63 | Malformed patterns `[`, `[]`, `[a-`, `***` | any | no panic | ✅ PASS |

### Parsing

| # | Input | Expected | Result |
|---|---|---|---|
| 64 | Empty file | 0 rules | ✅ PASS |
| 65 | `!linguist-vendored` (new) | "unspecified" → falls back | ✅ PASS |
| 66 | `linguist-language=a=b` | val="a=b" | ✅ PASS |
| 67 | `--linguist-vendored` | NOT parsed as linguist-vendored | ✅ PASS |
| 68 | CRLF line endings | rules parsed | ✅ PASS |
| 69 | Inline `# comment` | `#` as attr name | ✅ PASS (quirk) |
| 70 | `linguist-language=C++` / `C#` | val preserved | ✅ PASS |
| 71 | 1000 rules | all parsed | ✅ PASS |

### Attribute Semantics

| # | Scenario | Expected | Result |
|---|---|---|---|
| 72 | `-linguist-vendored` overrides default | false | ✅ PASS |
| 73 | Last rule wins: same attr | later value | ✅ PASS |
| 74 | `linguist-language=go` → alias resolved | "Go" | ✅ PASS |
| 75 | Unknown language → returned verbatim | "MyFakeLang9000" | ✅ PASS (quirk) |
| 76 | `IsDetectable` → (false, false) with no rules | fallback | ✅ PASS |
| 77 | `IsDetectable` → (true, true) with rule | forced in | ✅ PASS |
| 78 | `-linguist-detectable` → (false, true) | forced out | ✅ PASS |
| 79 | `!linguist-vendored` → unspecified → fallback | default | ✅ PASS |
| 80 | Multi-attr same line | both set | ✅ PASS |
| 81 | Real-world Kubernetes/TF patterns | correct | ✅ PASS |

### CLI Integration

| # | Scenario | Result |
|---|---|---|
| 82 | No `.gitattributes` | ✅ PASS |
| 83 | Empty `.gitattributes` | ✅ PASS |
| 84 | `**/vendor/**` does NOT match `notvendor/` | ✅ PASS (fixed) |
| 85 | `vendor/** -linguist-vendored` un-vendors files | ✅ PASS |
| 86 | `*.html linguist-language=Go` | ✅ PASS |
| 87 | `Vagrantfile -linguist-vendored` un-vendors it | ✅ PASS |
| 88 | 500 rules → < 50ms | ✅ PASS (36ms) |
| 89 | Conflicting rules → last wins | ✅ PASS |
| 90 | `vendor/ -linguist-vendored` fails for nested subdirs | ❌ FAIL (Bug #3) |

---

## Recommendation

### Must Fix Before Merge

**Bug #3: Trailing-slash pattern subdirectory matching**

The fix in commit `94aca8b` correctly handles files (`vendor/foo.go`) but not subdirectory paths as seen by `filepath.Walk` (`vendor/github.com/`). The CLI's directory skip logic means that any repository with nested dependencies under a `vendor/`-style subdirectory will have those files silently excluded from language statistics even when the user writes `vendor/ -linguist-vendored`.

The workaround is to use `vendor/**` instead of `vendor/`. This works correctly and is well-tested. However, real Linguist `.gitattributes` files commonly use the trailing-slash form, so this bug will be a source of user confusion.

**Fix:** In `matchGitPattern`, when the path ends with `/` inside the trailing-slash branch, check if the stripped path is the directory itself or a child:

```go
if strings.HasSuffix(path, "/") {
    dirPath := strings.TrimSuffix(path, "/")
    if dirPath == dir || strings.HasPrefix(dirPath, dir+"/") {
        return true
    }
    return false
}
```

### Still Applies from Round 1

- Document `IsDetectable` asymmetry (now returning `(bool, bool)` — documented in godoc ✅)
- Warn on unknown language names in `GetLanguage`
- Document inline-comment limitation
- (Future) Escaped spaces in filenames
- (Future) Multi-level `.gitattributes` files
- (Future) Case-insensitive matching on Windows/macOS

---

## Appendix: Reproduce Bug #3

```bash
mkdir /tmp/bug3-test && cd /tmp/bug3-test
mkdir -p vendor/github.com/pkg
echo 'package p' > vendor/github.com/pkg/errors.go
cat > .gitattributes << 'EOF'
vendor/ -linguist-vendored
EOF
enry -breakdown .
# Expected: vendor/github.com/pkg/errors.go appears (un-vendored)
# Actual:   absent — vendor/github.com/ is still IsVendor=true, causing SkipDir

# Workaround:
cat > .gitattributes << 'EOF'
vendor/** -linguist-vendored
EOF
enry -breakdown .
# Now errors.go correctly appears
```
