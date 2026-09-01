# Gitignore & Gitattributes Pattern Matching Specification

> Extracted from Git source (commit cd412a4962). Intended as a self-contained reference
> for implementing `.gitignore` and `.gitattributes` pattern matching in Go.
> All concepts are explained inline — no Git internals knowledge required.
>
> Source files: [`wildmatch.c`](https://github.com/git/git/blob/master/wildmatch.c), [`wildmatch.h`](https://github.com/git/git/blob/master/wildmatch.h), [`dir.c`](https://github.com/git/git/blob/master/dir.c), [`dir.h`](https://github.com/git/git/blob/master/dir.h), [`attr.c`](https://github.com/git/git/blob/master/attr.c), [`attr.h`](https://github.com/git/git/blob/master/attr.h),
> [`t/t3070-wildmatch.sh`](https://github.com/git/git/blob/master/t/t3070-wildmatch.sh), [`t/t0008-ignores.sh`](https://github.com/git/git/blob/master/t/t0008-ignores.sh), [`t/t3001-ls-files-others-exclude.sh`](https://github.com/git/git/blob/master/t/t3001-ls-files-others-exclude.sh), [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh),
> [`Documentation/gitignore.adoc`](https://github.com/git/git/blob/master/Documentation/gitignore.adoc), [`Documentation/gitattributes.adoc`](https://github.com/git/git/blob/master/Documentation/gitattributes.adoc)

---

## Part 1: Wildmatch — The Pattern Matching Engine

### 1.1 What Is Wildmatch?

Wildmatch is Git's glob-style pattern matching function. It takes a **pattern** and a **text**
string and returns whether the pattern matches the text. It is the engine behind all
`.gitignore` pattern matching, ref filtering, pathspec matching, etc.

```
func wildmatch(pattern, text string, flags uint) bool
```

### 1.2 Matching Modes

Wildmatch has two independent boolean flags that combine into 4 modes:

| Flag | Name | Effect |
|------|------|--------|
| `WM_PATHNAME` | **Pathname mode** | The `/` character is treated as a special directory separator. `*`, `?`, and `[…]` **cannot match** `/`. Only `**` can cross `/` boundaries. |
| `WM_CASEFOLD` | **Case-insensitive** | Matching is case-insensitive. `A` matches `a`. Character ranges like `[A-Z]` also match lowercase. POSIX classes like `[:upper:]` match lowercase too. |

The 4 combinations:

| Column abbrev | Flags | Used by |
|---------------|-------|---------|
| **PN** | `WM_PATHNAME` | `.gitignore` matching against full relative paths |
| **PN+CF** | `WM_PATHNAME \| WM_CASEFOLD` | `.gitignore` on case-insensitive filesystems |
| **PM** | *(none)* | `.gitignore` matching against basenames only (no `/` in pattern) |
| **PM+CF** | `WM_CASEFOLD` | basename matching, case-insensitive |

**For `.gitignore`**: patterns containing `/` are matched with `WM_PATHNAME` against the full
relative path. Patterns without `/` are matched with no flags against the basename only.
Case-insensitive mode (`WM_CASEFOLD`) is added when the filesystem is case-insensitive
(e.g., macOS HFS+, Windows NTFS).

### 1.3 Pattern Syntax

| Syntax | Meaning |
|--------|---------|
| `*` | Matches any sequence of characters **except** `/` (in pathname mode). Without pathname mode, matches anything. |
| `**` | In pathname mode and when at a path boundary (start of pattern, after `/`, or before `/`/end): matches any sequence including `/` — i.e., zero or more directory levels. Otherwise treated as `*`. |
| `?` | Matches exactly one character **except** `/` (in pathname mode). Without pathname mode, matches any character. |
| `[abc]` | Character class — matches one of the listed characters. |
| `[a-z]` | Character range — matches one character in the range (inclusive). |
| `[!abc]` or `[^abc]` | Negated character class — matches one character NOT in the set. In pathname mode, `[…]` never matches `/` regardless of contents. |
| `[[:alpha:]]` | POSIX character class. Supported: `alnum`, `alpha`, `blank`, `cntrl`, `digit`, `graph`, `lower`, `print`, `punct`, `space`, `upper`, `xdigit`. |
| `\x` | Escape — the next character is matched literally (e.g., `\*` matches a literal `*`, `\\` matches `\`, `\a` matches `a`). |

**`**` boundary rules** (pathname mode only):
- `**/foo` — `**` at start followed by `/`: matches `foo` at any depth
- `foo/**` — `**` at end preceded by `/`: matches everything inside `foo/`
- `foo/**/bar` — `/**/` in middle: matches `foo/bar`, `foo/x/bar`, `foo/x/y/bar`, etc. (zero or more directories)
- `foo**bar` — `**` NOT at a boundary: treated as regular `*` (cannot cross `/`)

### 1.4 Return Values

- **1** (match): pattern matches the text
- **0** (no match): pattern does not match

Internally, the algorithm uses abort codes to prevent exponential backtracking, but these
are not exposed to callers.

---

### 1.5 Test Cases — Basic Features

These test cases exercise literals, `?`, `*`, `\` escaping, and character classes.

In all tables: **1** = match, **0** = no match.

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 1 | `foo` | `foo` | 1 | 1 | 1 | 1 | Exact literal match |
| 2 | `foo` | `bar` | 0 | 0 | 0 | 0 | Different literal |
| 3 | *(empty)* | *(empty)* | 1 | 1 | 1 | 1 | Both empty strings match |
| 4 | `foo` | `???` | 1 | 1 | 1 | 1 | Three `?` match three chars |
| 5 | `foo` | `??` | 0 | 0 | 0 | 0 | Two `?` don't match three chars |
| 6 | `foo` | `*` | 1 | 1 | 1 | 1 | `*` matches everything |
| 7 | `foo` | `f*` | 1 | 1 | 1 | 1 | `f*` matches `foo` |
| 8 | `foo` | `*f` | 0 | 0 | 0 | 0 | `*f` — `foo` doesn't end with `f` |
| 9 | `foo` | `*foo*` | 1 | 1 | 1 | 1 | `*foo*` matches |
| 10 | `foobar` | `*ob*a*r*` | 1 | 1 | 1 | 1 | Multiple `*` segments |
| 11 | `aaaaaaabababab` | `*ab` | 1 | 1 | 1 | 1 | Repetitive text, ends with `ab` |
| 12 | `foo*` | `foo\*` | 1 | 1 | 1 | 1 | `\*` matches literal `*` in text |
| 13 | `foobar` | `foo\*bar` | 0 | 0 | 0 | 0 | `\*` requires literal `*` between foo and bar |
| 14 | `f\oo` | `f\\oo` | 1 | 1 | 1 | 1 | `\\` matches literal `\` |
| 15 | `foo\` | `foo\` | 0 | 0 | 0 | 0 | Trailing `\` in pattern is invalid — no match |
| 16 | `ball` | `*[al]?` | 1 | 1 | 1 | 1 | `*` then char class then `?` |
| 17 | `ten` | `[ten]` | 0 | 0 | 0 | 0 | `[ten]` matches ONE char, not three |
| 18 | `ten` | `**[!te]` | 1 | 1 | 1 | 1 | `**` matches `te`, `[!te]` matches `n` |
| 19 | `ten` | `**[!ten]` | 0 | 0 | 0 | 0 | All chars of `ten` are in `[ten]` |
| 20 | `ten` | `t[a-g]n` | 1 | 1 | 1 | 1 | `e` is in range `a-g` |
| 21 | `ten` | `t[!a-g]n` | 0 | 0 | 0 | 0 | `e` IS in `a-g`, negated → no match |
| 22 | `ton` | `t[!a-g]n` | 1 | 1 | 1 | 1 | `o` is NOT in `a-g`, negated → match |
| 23 | `ton` | `t[^a-g]n` | 1 | 1 | 1 | 1 | `^` is alias for `!` in character classes |
| 24 | `a]b` | `a[]]b` | 1 | 1 | 1 | 1 | `]` immediately after `[` is literal |
| 25 | `a-b` | `a[]-]b` | 1 | 1 | 1 | 1 | Class `[]-]` contains `]` and `-` |
| 26 | `a]b` | `a[]-]b` | 1 | 1 | 1 | 1 | `]` matches `[]-]` |
| 27 | `aab` | `a[]-]b` | 0 | 0 | 0 | 0 | `a` is not `]` or `-` |
| 28 | `aab` | `a[]a-]b` | 1 | 1 | 1 | 1 | Class `[]a-]` contains `]`, `a`, `-` |
| 29 | `]` | `]` | 1 | 1 | 1 | 1 | Literal `]` outside brackets |

### 1.6 Test Cases — Slash / Pathname Matching

These test the interaction between `*`, `**`, `?`, `[…]` and the `/` separator.
The key insight: in **pathname mode** (PN), `/` is special and `*`/`?`/`[…]` cannot cross it.
In **plain mode** (PM), `/` is an ordinary character.

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 30 | `foo/baz/bar` | `foo*bar` | 0 | 0 | 1 | 1 | `*` can't cross `/` in PN |
| 31 | `foo/baz/bar` | `foo**bar` | 0 | 0 | 1 | 1 | `**` not at boundary → treated as `*` in PN |
| 32 | `foobazbar` | `foo**bar` | 1 | 1 | 1 | 1 | No `/` in text → `**` works like `*` |
| 33 | `foo/baz/bar` | `foo/**/bar` | 1 | 1 | 1 | 1 | `/**/` matches one dir level |
| 34 | `foo/baz/bar` | `foo/**/**/bar` | 1 | 1 | 0 | 0 | Multiple `/**/` works in PN; PM treats `**` literally |
| 35 | `foo/b/a/z/bar` | `foo/**/bar` | 1 | 1 | 1 | 1 | `/**/` matches multiple dir levels |
| 36 | `foo/b/a/z/bar` | `foo/**/**/bar` | 1 | 1 | 1 | 1 | Multiple `/**/` matches deep paths |
| 37 | `foo/bar` | `foo/**/bar` | 1 | 1 | 0 | 0 | `/**/` matches **zero** dirs in PN |
| 38 | `foo/bar` | `foo/**/**/bar` | 1 | 1 | 0 | 0 | Multiple `/**/` also matches zero dirs in PN |
| 39 | `foo/bar` | `foo?bar` | 0 | 0 | 1 | 1 | `?` can't match `/` in PN |
| 40 | `foo/bar` | `foo[/]bar` | 0 | 0 | 1 | 1 | `[/]` can't match `/` in PN |
| 41 | `foo/bar` | `foo[^a-z]bar` | 0 | 0 | 1 | 1 | Negated class can't match `/` in PN |
| 42 | `foo/bar` | `f[^eiu][^eiu][^eiu][^eiu][^eiu]r` | 0 | 0 | 1 | 1 | Multiple negated classes can't cross `/` in PN |
| 43 | `foo-bar` | `f[^eiu][^eiu][^eiu][^eiu][^eiu]r` | 1 | 1 | 1 | 1 | No `/` in text → works fine |
| 44 | `foo` | `**/foo` | 1 | 1 | 0 | 0 | `**/` at start matches at root in PN; PM: literal `**/` |
| 45 | `XXX/foo` | `**/foo` | 1 | 1 | 1 | 1 | `**/foo` matches at any depth |
| 46 | `bar/baz/foo` | `**/foo` | 1 | 1 | 1 | 1 | `**/foo` matches deep |
| 47 | `bar/baz/foo` | `*/foo` | 0 | 0 | 1 | 1 | `*/foo` only matches one dir in PN |
| 48 | `foo/bar/baz` | `**/bar*` | 0 | 0 | 1 | 1 | `**/bar*` — `**` at boundary but `bar*` can't match `bar/baz` in PN |
| 49 | `deep/foo/bar/baz` | `**/bar/*` | 1 | 1 | 1 | 1 | `**/bar/*` matches |
| 50 | `deep/foo/bar/baz/` | `**/bar/*` | 0 | 0 | 1 | 1 | Trailing `/` in text: `*` can't match `baz/` in PN |
| 51 | `deep/foo/bar/baz/` | `**/bar/**` | 1 | 1 | 1 | 1 | `/**` at end matches trailing `/` |
| 52 | `deep/foo/bar` | `**/bar/*` | 0 | 0 | 0 | 0 | Nothing after `bar/` for `*` to match |
| 53 | `deep/foo/bar/` | `**/bar/**` | 1 | 1 | 1 | 1 | `/**` matches the trailing `/` |
| 54 | `foo/bar/baz` | `**/bar**` | 0 | 0 | 1 | 1 | `bar**` — `**` not at boundary in PN |
| 55 | `foo/bar/baz/x` | `*/bar/**` | 1 | 1 | 1 | 1 | `*/` matches one dir, `/**` matches rest |
| 56 | `deep/foo/bar/baz/x` | `*/bar/**` | 0 | 0 | 1 | 1 | `*/` only matches one dir in PN |
| 57 | `deep/foo/bar/baz/x` | `**/bar/*/*` | 1 | 1 | 1 | 1 | `**/bar` then `/*` twice |

### 1.7 Test Cases — Various Additional

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 58 | `acrt` | `a[c-c]st` | 0 | 0 | 0 | 0 | `[c-c]` is just `c`; `acrt` ≠ `acst` |
| 59 | `acrt` | `a[c-c]rt` | 1 | 1 | 1 | 1 | `[c-c]` matches `c`, rest matches |
| 60 | `]` | `[!]-]` | 0 | 0 | 0 | 0 | `]` IS in `]-]`, negated → no match |
| 61 | `a` | `[!]-]` | 1 | 1 | 1 | 1 | `a` is NOT in `]-]`, negated → match |
| 62 | *(empty)* | `\` | 0 | 0 | 0 | 0 | Lone `\` is invalid pattern |
| 63 | `\` | `\` | 0 | 0 | 0 | 0 | Lone `\` pattern: invalid, never matches |
| 64 | `XXX/\` | `*/\` | 0 | 0 | 0 | 0 | Trailing `\` in pattern is invalid |
| 65 | `XXX/\` | `*/\\` | 1 | 1 | 1 | 1 | `\\` matches literal `\` |
| 66 | `foo` | `foo` | 1 | 1 | 1 | 1 | Exact match (repeated for context) |
| 67 | `@foo` | `@foo` | 1 | 1 | 1 | 1 | `@` has no special meaning |
| 68 | `foo` | `@foo` | 0 | 0 | 0 | 0 | `@` is literal |
| 69 | `[ab]` | `\[ab]` | 1 | 1 | 1 | 1 | `\[` matches literal `[` |
| 70 | `[ab]` | `[[]ab]` | 1 | 1 | 1 | 1 | `[[]` is char class containing `[` |
| 71 | `[ab]` | `[[:]ab]` | 1 | 1 | 1 | 1 | `[[:` without closing `:]` — treated as chars `[`, `:` |
| 72 | `[ab]` | `[[::]ab]` | 0 | 0 | 0 | 0 | `[::]` is invalid POSIX class → abort |
| 73 | `[ab]` | `[[:digit]ab]` | 1 | 1 | 1 | 1 | `[:digit]` missing final `:` → treated as chars |
| 74 | `[ab]` | `[\[:]ab]` | 1 | 1 | 1 | 1 | `\[` in class is literal `[` |
| 75 | `?a?b` | `\??\?b` | 1 | 1 | 1 | 1 | `\?` matches literal `?`, `?` matches `a` |
| 76 | `abc` | `\a\b\c` | 1 | 1 | 1 | 1 | `\a` matches `a` (escape of non-special char) |
| 77 | `foo` | *(empty)* | 0 | 0 | 0 | 0 | Empty pattern does not match non-empty text |
| 78 | `foo/bar/baz/to` | `**/t[o]` | 1 | 1 | 1 | 1 | `**/` then char class `[o]` |

### 1.8 Test Cases — POSIX Character Classes

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 79 | `a1B` | `[[:alpha:]][[:digit:]][[:upper:]]` | 1 | 1 | 1 | 1 | `a`=alpha, `1`=digit, `B`=upper |
| 80 | `a` | `[[:digit:][:upper:][:space:]]` | 0 | 1 | 0 | 1 | `a` is not digit/upper/space; with casefold, `[:upper:]` matches lowercase |
| 81 | `A` | `[[:digit:][:upper:][:space:]]` | 1 | 1 | 1 | 1 | `A` is `[:upper:]` |
| 82 | `1` | `[[:digit:][:upper:][:space:]]` | 1 | 1 | 1 | 1 | `1` is `[:digit:]` |
| 83 | `1` | `[[:digit:][:upper:][:spaci:]]` | 0 | 0 | 0 | 0 | `[:spaci:]` is invalid → abort, no match |
| 84 | ` ` | `[[:digit:][:upper:][:space:]]` | 1 | 1 | 1 | 1 | space is `[:space:]` |
| 85 | `.` | `[[:digit:][:upper:][:space:]]` | 0 | 0 | 0 | 0 | `.` is none of those |
| 86 | `.` | `[[:digit:][:punct:][:space:]]` | 1 | 1 | 1 | 1 | `.` is `[:punct:]` |
| 87 | `5` | `[[:xdigit:]]` | 1 | 1 | 1 | 1 | `5` is hex digit |
| 88 | `f` | `[[:xdigit:]]` | 1 | 1 | 1 | 1 | `f` is hex digit |
| 89 | `D` | `[[:xdigit:]]` | 1 | 1 | 1 | 1 | `D` is hex digit |
| 90 | `_` | `[[:alnum:][:alpha:][:blank:][:cntrl:][:digit:][:graph:][:lower:][:print:][:punct:][:space:][:upper:][:xdigit:]]` | 1 | 1 | 1 | 1 | `_` matches `[:punct:]` |
| 91 | `.` | `[^[:alnum:][:alpha:][:blank:][:cntrl:][:digit:][:lower:][:space:][:upper:][:xdigit:]]` | 1 | 1 | 1 | 1 | `.` is not in any of those (but is `[:graph:]`/`[:print:]`/`[:punct:]`, which are excluded from the negated set) |
| 92 | `5` | `[a-c[:digit:]x-z]` | 1 | 1 | 1 | 1 | `5` matches `[:digit:]` |
| 93 | `b` | `[a-c[:digit:]x-z]` | 1 | 1 | 1 | 1 | `b` matches `a-c` |
| 94 | `y` | `[a-c[:digit:]x-z]` | 1 | 1 | 1 | 1 | `y` matches `x-z` |
| 95 | `q` | `[a-c[:digit:]x-z]` | 0 | 0 | 0 | 0 | `q` is not in any range/class |

### 1.9 Test Cases — Malformed / Edge-Case Patterns

These test bracket expressions with unusual or incomplete syntax.

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 96 | `]` | `[\-^]` | 1 | 1 | 1 | 1 | `\-` is literal `-`; class is `-` and `^`; but `]` matches `^`? No — range `-` to `^` (ASCII 45-94) includes `]` (93) |
| 97 | `[` | `[\-^]` | 0 | 0 | 0 | 0 | `[` is ASCII 91, within range `-`(45) to `^`(94)… but `\-` escapes `-` as literal, so class has escaped `-` then range check with `^`. Actually: `\-` is literal `-`, then `-^` is not a range because `-` was consumed by escape. So class is {`-`, `^`}. `[` matches neither. |
| 98 | `-` | `[\-_]` | 1 | 1 | 1 | 1 | `\-` is literal `-`; then range `-` to `_`. `-` matches. |
| 99 | `]` | `[\]]` | 1 | 1 | 1 | 1 | `\]` is literal `]` in class |
| 100 | `\]` | `[\]]` | 0 | 0 | 0 | 0 | Class matches single char `]`, not `\]` |
| 101 | `\` | `[\]]` | 0 | 0 | 0 | 0 | `\` is not `]` |
| 102 | `ab` | `a[]b` | 0 | 0 | 0 | 0 | Empty bracket `[]` — never completes, no match |
| 103 | `a[]b` | `a[]b` | 0 | 0 | 0 | 0 | Wildmatch: incomplete bracket → no match |
| 104 | `ab[` | `ab[` | 0 | 0 | 0 | 0 | Wildmatch: incomplete bracket → no match |
| 105 | `ab` | `[!` | 0 | 0 | 0 | 0 | Incomplete negated bracket |
| 106 | `ab` | `[-` | 0 | 0 | 0 | 0 | Incomplete bracket with `-` |
| 107 | `-` | `[-]` | 1 | 1 | 1 | 1 | `-` at start of class is literal |
| 108 | `-` | `[a-` | 0 | 0 | 0 | 0 | Incomplete range |
| 109 | `-` | `[!a-` | 0 | 0 | 0 | 0 | Incomplete negated range |
| 110 | `-` | `[--A]` | 1 | 1 | 1 | 1 | Range `-` (45) to `A` (65); `-` is in range |
| 111 | `5` | `[--A]` | 1 | 1 | 1 | 1 | `5` (53) is in range 45-65 |
| 112 | ` ` | `[ --]` | 1 | 1 | 1 | 1 | Range ` ` (32) to `-` (45); space (32) is in range |
| 113 | `$` | `[ --]` | 1 | 1 | 1 | 1 | `$` (36) is in range 32-45 |
| 114 | `-` | `[ --]` | 1 | 1 | 1 | 1 | `-` (45) is in range 32-45 |
| 115 | `0` | `[ --]` | 0 | 0 | 0 | 0 | `0` (48) is NOT in range 32-45 |
| 116 | `-` | `[---]` | 1 | 1 | 1 | 1 | Range `-` to `-`; matches `-` |
| 117 | `-` | `[------]` | 1 | 1 | 1 | 1 | Multiple `-` sequences; matches `-` |
| 118 | `j` | `[a-e-n]` | 0 | 0 | 0 | 0 | Range `a-e` then literal `-` then `n`; `j` matches none |
| 119 | `-` | `[a-e-n]` | 1 | 1 | 1 | 1 | `-` matches the literal `-` after the range |
| 120 | `a` | `[!------]` | 1 | 1 | 1 | 1 | Negated class containing only `-`; `a` is not `-` |
| 121 | `[` | `[]-a]` | 0 | 0 | 0 | 0 | Range `]` (93) to `a` (97); `[` (91) is not in range |
| 122 | `^` | `[]-a]` | 1 | 1 | 1 | 1 | `^` (94) is in range 93-97 |
| 123 | `^` | `[!]-a]` | 0 | 0 | 0 | 0 | Negated: `^` IS in range → no match |
| 124 | `[` | `[!]-a]` | 1 | 1 | 1 | 1 | Negated: `[` is NOT in range → match |
| 125 | `^` | `[a^bc]` | 1 | 1 | 1 | 1 | `^` as literal char in class (not at start) |
| 126 | `-b]` | `[a-]b]` | 1 | 1 | 1 | 1 | `[a-]` is class with `a` and `-` (before `]`); then literal `b]` |
| 127 | `\` | `[\]` | 0 | 0 | 0 | 0 | `\]` escapes `]`, bracket never closes → no match |
| 128 | `\` | `[\\]` | 1 | 1 | 1 | 1 | `\\` in class is literal `\` |
| 129 | `\` | `[!\\]` | 0 | 0 | 0 | 0 | Negated: `\` IS `\\` → no match |
| 130 | `G` | `[A-\\]` | 1 | 1 | 1 | 1 | Range `A` (65) to `\` (92); `G` (71) is in range |
| 131 | `aaabbb` | `b*a` | 0 | 0 | 0 | 0 | Text doesn't start with `b` |
| 132 | `aabcaa` | `*ba*` | 0 | 0 | 0 | 0 | No `ba` substring in text |
| 133 | `,` | `[,]` | 1 | 1 | 1 | 1 | Literal `,` in class |
| 134 | `,` | `[\\,]` | 1 | 1 | 1 | 1 | Class: `\\` (literal `\`) and `,`; `,` matches |
| 135 | `\` | `[\\,]` | 1 | 1 | 1 | 1 | `\` matches `\\` in class |
| 136 | `-` | `[,-.]` | 1 | 1 | 1 | 1 | Range `,` (44) to `.` (46); `-` (45) is in range |
| 137 | `+` | `[,-.]` | 0 | 0 | 0 | 0 | `+` (43) is NOT in range 44-46 |
| 138 | `-.]` | `[,-.]` | 0 | 0 | 0 | 0 | Class matches single char, text is 3 chars |
| 139 | `2` | `[\1-\3]` | 1 | 1 | 1 | 1 | `\1` to `\3` (escaped digits = literal `1`-`3`); `2` in range |
| 140 | `3` | `[\1-\3]` | 1 | 1 | 1 | 1 | `3` in range |
| 141 | `4` | `[\1-\3]` | 0 | 0 | 0 | 0 | `4` not in range `1`-`3` |
| 142 | `\` | `[[-\]]` | 1 | 1 | 1 | 1 | Range `[` (91) to `\]` (escaped `]` = 93); `\` (92) in range |
| 143 | `[` | `[[-\]]` | 1 | 1 | 1 | 1 | `[` (91) in range 91-93 |
| 144 | `]` | `[[-\]]` | 1 | 1 | 1 | 1 | `]` (93) in range 91-93 |
| 145 | `-` | `[[-\]]` | 0 | 0 | 0 | 0 | `-` (45) NOT in range 91-93 |

### 1.10 Test Cases — Recursion / Complex Multi-Wildcard

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 146 | `-adobe-courier-bold-o-normal--12-120-75-75-m-70-iso8859-1` | `-*-*-*-*-*-*-12-*-*-*-m-*-*-*` | 1 | 1 | 1 | 1 | Complex multi-`*` pattern |
| 147 | `-adobe-courier-bold-o-normal--12-120-75-75-X-70-iso8859-1` | `-*-*-*-*-*-*-12-*-*-*-m-*-*-*` | 0 | 0 | 0 | 0 | `X` where `m` is expected |
| 148 | `-adobe-courier-bold-o-normal--12-120-75-75-/-70-iso8859-1` | `-*-*-*-*-*-*-12-*-*-*-m-*-*-*` | 0 | 0 | 0 | 0 | `/` in text; even in PM, `m` is required not `/` |
| 149 | `XXX/adobe/courier/bold/o/normal//12/120/75/75/m/70/iso8859/1` | `XXX/*/*/*/*/*/*/12/*/*/*/m/*/*/*` | 1 | 1 | 1 | 1 | Path with `//` (empty segment) |
| 150 | `XXX/adobe/courier/bold/o/normal//12/120/75/75/X/70/iso8859/1` | `XXX/*/*/*/*/*/*/12/*/*/*/m/*/*/*` | 0 | 0 | 0 | 0 | `X` where `m` expected |
| 151 | `abcd/abcdefg/abcdefghijk/abcdefghijklmnop.txt` | `**/*a*b*g*n*t` | 1 | 1 | 1 | 1 | Deep path, complex suffix |
| 152 | `abcd/abcdefg/abcdefghijk/abcdefghijklmnop.txtz` | `**/*a*b*g*n*t` | 0 | 0 | 0 | 0 | `txtz` doesn't end with `t` |
| 153 | `foo` | `*/*/*` | 0 | 0 | 0 | 0 | 0 slashes, need 2 |
| 154 | `foo/bar` | `*/*/*` | 0 | 0 | 0 | 0 | 1 slash, need 2 |
| 155 | `foo/bba/arr` | `*/*/*` | 1 | 1 | 1 | 1 | Exactly 2 slashes → 3 segments |
| 156 | `foo/bb/aa/rr` | `*/*/*` | 0 | 0 | 1 | 1 | 3 slashes: PN `*` can't cross `/`; PM can |
| 157 | `foo/bb/aa/rr` | `**/**/**` | 1 | 1 | 1 | 1 | `**` crosses all `/` boundaries |
| 158 | `abcXdefXghi` | `*X*i` | 1 | 1 | 1 | 1 | No `/` → all modes agree |
| 159 | `ab/cXd/efXg/hi` | `*X*i` | 0 | 0 | 1 | 1 | PN: `*` can't cross `/` |
| 160 | `ab/cXd/efXg/hi` | `*/*X*/*/*i` | 1 | 1 | 1 | 1 | Explicit `/` in pattern matches structure |
| 161 | `ab/cXd/efXg/hi` | `**/*X*/**/*i` | 1 | 1 | 1 | 1 | `**` and `*` work together |

### 1.11 Test Cases — Extra Pathmatch

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 162 | `foo` | `fo` | 0 | 0 | 0 | 0 | Incomplete pattern |
| 163 | `foo/bar` | `foo/bar` | 1 | 1 | 1 | 1 | Exact path |
| 164 | `foo/bar` | `foo/*` | 1 | 1 | 1 | 1 | `*` matches `bar` |
| 165 | `foo/bba/arr` | `foo/*` | 0 | 0 | 1 | 1 | PN: `*` can't match `bba/arr` |
| 166 | `foo/bba/arr` | `foo/**` | 1 | 1 | 1 | 1 | `/**` matches everything below |
| 167 | `foo/bba/arr` | `foo*` | 0 | 0 | 1 | 1 | PN: `*` after `foo` can't cross `/` |
| 168 | `foo/bba/arr` | `foo**` | 0 | 0 | 1 | 1 | PN: `**` not at boundary (no preceding `/`) → treated as `*` |
| 169 | `foo/bba/arr` | `foo/*arr` | 0 | 0 | 1 | 1 | PN: `*` matches only within one dir |
| 170 | `foo/bba/arr` | `foo/**arr` | 0 | 0 | 1 | 1 | PN: `**` before `arr` not at boundary → like `*` |
| 171 | `foo/bba/arr` | `foo/*z` | 0 | 0 | 0 | 0 | No `z` anywhere |
| 172 | `foo/bba/arr` | `foo/**z` | 0 | 0 | 0 | 0 | No `z` anywhere |
| 173 | `foo/bar` | `foo?bar` | 0 | 0 | 1 | 1 | PN: `?` can't match `/` |
| 174 | `foo/bar` | `foo[/]bar` | 0 | 0 | 1 | 1 | PN: `[/]` can't match `/` |
| 175 | `foo/bar` | `foo[^a-z]bar` | 0 | 0 | 1 | 1 | PN: negated class can't match `/` |
| 176 | `ab/cXd/efXg/hi` | `*Xg*i` | 0 | 0 | 1 | 1 | PN: `*` can't cross multiple `/` |

### 1.12 Test Cases — Case Sensitivity

These show how `WM_CASEFOLD` affects character ranges and POSIX classes.

| # | Text | Pattern | PN | PN+CF | PM | PM+CF | Note |
|---|------|---------|:--:|:-----:|:--:|:-----:|------|
| 177 | `a` | `[A-Z]` | 0 | 1 | 0 | 1 | `a` not in `A-Z`; casefold: `A` in `A-Z` |
| 178 | `A` | `[A-Z]` | 1 | 1 | 1 | 1 | `A` in `A-Z` |
| 179 | `A` | `[a-z]` | 0 | 1 | 0 | 1 | `A` not in `a-z`; casefold: `a` in `a-z` |
| 180 | `a` | `[a-z]` | 1 | 1 | 1 | 1 | `a` in `a-z` |
| 181 | `a` | `[[:upper:]]` | 0 | 1 | 0 | 1 | `a` is not upper; casefold: `[:upper:]` matches lowercase too |
| 182 | `A` | `[[:upper:]]` | 1 | 1 | 1 | 1 | `A` is upper |
| 183 | `A` | `[[:lower:]]` | 0 | 1 | 0 | 1 | `A` is not lower; casefold matches |
| 184 | `a` | `[[:lower:]]` | 1 | 1 | 1 | 1 | `a` is lower |
| 185 | `A` | `[B-Za]` | 0 | 1 | 0 | 1 | `A` not in `B-Z` or `a`; casefold: `a` matches literal `a` in class |
| 186 | `a` | `[B-Za]` | 1 | 1 | 1 | 1 | `a` matches literal `a` in class |
| 187 | `A` | `[B-a]` | 0 | 1 | 0 | 1 | Range `B`(66) to `a`(97); `A`(65) not in range; casefold matches |
| 188 | `a` | `[B-a]` | 1 | 1 | 1 | 1 | `a`(97) is in range 66-97 |
| 189 | `z` | `[Z-y]` | 0 | 1 | 0 | 1 | Range `Z`(90) to `y`(121); `z`(122) not in range; casefold: `Z`(90) is in range |
| 190 | `Z` | `[Z-y]` | 1 | 1 | 1 | 1 | `Z`(90) is in range 90-121 |

### 1.13 Test Case — Exponential Backtracking Prevention

This test ensures the matching engine does **not** exhibit exponential time complexity
on pathological patterns. An implementation must complete (or abort) in bounded time.

| Text | Pattern | Expected |
|------|---------|----------|
| `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab` (59 `a`s + `b`) | `*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a` (16 `*a` segments) | **No match**, must complete in < 2 seconds |

A naive recursive implementation would try O(2^n) paths. The correct implementation
uses abort codes to prune the search space.

---

## Part 2: Gitignore Pattern Parsing

This section describes how lines in a `.gitignore` file are parsed into match rules.

### 2.1 File Format

A `.gitignore` file is a plain text file with one pattern per line. The parser processes it
line by line with these rules (applied in order):

1. **UTF-8 BOM**: If the file starts with a UTF-8 BOM (`\xEF\xBB\xBF`), skip it.
2. **Line endings**: Both `\n` (Unix) and `\r\n` (Windows) are recognized as line terminators.
3. **Blank lines**: Lines containing only whitespace are ignored. They serve as visual separators.
4. **Comment lines**: Lines where the first character is `#` are comments and ignored.
   To match a file starting with `#`, escape it: `\#`.
5. **Trailing whitespace**: Trailing spaces are stripped from the pattern, **unless** they
   are escaped with `\`. For example:
   - `foo   ` → parsed as `foo` (trailing spaces stripped)
   - `foo\ ` → parsed as `foo ` (one trailing space preserved)
   - `foo\\ ` → parsed as `foo\` (the `\\` is a literal `\`, trailing space stripped)
6. **Negation prefix**: If the line starts with `!`, it's a negation pattern — it
   **un-ignores** files that were previously ignored by an earlier pattern.
   To match a file starting with `!`, escape it: `\!`.
7. **Trailing `/`**: If the pattern ends with `/`, the slash is removed and the pattern
   is marked as **directory-only** — it will only match directories, not regular files.
8. **Leading `/`**: A leading `/` anchors the pattern to the directory containing the
   `.gitignore` file (see anchoring rules below). The `/` itself is stripped from the pattern.

### 2.2 Pattern Flags

After parsing, each pattern has these computed flags:

| Flag | Condition | Effect |
|------|-----------|--------|
| **Negative** | Line started with `!` | Match means "un-ignore" instead of "ignore" |
| **MustBeDir** | Line ended with `/` | Pattern only matches directories |
| **NoDir** | Pattern contains no `/` (after stripping leading `/` and trailing `/`) | Pattern is matched against the **basename** only (not the full path) |
| **EndsWith** | Pattern starts with `*` and the rest has no wildcards (e.g., `*.txt`) | Optimization: can use suffix matching instead of full wildmatch |

### 2.3 Anchoring Rules

This is one of the most important and subtle aspects of `.gitignore`:

- **Unanchored** (no `/` in pattern): The pattern matches the **filename (basename) only**,
  at any directory depth. Example: `*.o` ignores `foo.o`, `src/bar.o`, `deep/nested/baz.o`.

- **Anchored** (pattern contains `/` in beginning or middle): The pattern is matched against
  the **full path relative to the `.gitignore` location**, using pathname mode (where `*` can't
  cross `/`). Example: `doc/frotz` only matches `doc/frotz` relative to the `.gitignore`, not
  `a/doc/frotz`.

- A **leading `/`** forces anchoring but is then stripped. So `/foo` is equivalent to an anchored
  pattern `foo` — it matches `foo` in the `.gitignore`'s directory only, not `a/foo`.

- A **trailing `/`** does NOT count as a `/` for anchoring purposes. So `build/` is treated as
  unanchored (the `/` is stripped and the pattern has no remaining `/`). It matches `build` and
  `a/build` directories at any depth, but NOT regular files named `build`.

- Both `doc/frotz` and `/doc/frotz` behave identically — a leading slash is redundant when
  there is already a middle slash.

### 2.4 Matching Algorithm for a Single `.gitignore`

Given a list of parsed patterns and a path to check:

```
for each pattern in REVERSE order (last pattern first):
    if pattern.MustBeDir and path is not a directory:
        skip this pattern

    if pattern.NoDir:
        # Pattern has no slash → match against basename only
        result = wildmatch(pattern, basename(path), flags=0)
    else:
        # Pattern has slash → match against full relative path
        result = wildmatch(pattern, relativePath, flags=WM_PATHNAME)

    if result == match:
        if pattern.Negative:
            return NOT_IGNORED
        else:
            return IGNORED

return UNDECIDED  # no pattern matched
```

Key points:
- **Last match wins**: patterns are checked from bottom to top; the first match found
  (= the last matching line in the file) determines the outcome.
- **Basename matching uses no flags** (plain mode): this means `*` CAN match `/` in the
  basename... but basenames don't contain `/`, so this is only relevant for unanchored
  patterns.
- **Full path matching uses `WM_PATHNAME`**: `*` cannot cross `/`, `**` can.
- If `WM_CASEFOLD` is needed (case-insensitive filesystem), add it to the flags.

### 2.5 Parsing Test Cases

These test the line-parsing rules from section 2.1:

| Input line | Parsed pattern | Negative | MustBeDir | NoDir | Note |
|------------|---------------|:--------:|:---------:|:-----:|------|
| `foo` | `foo` | | | yes | Simple basename pattern |
| `*.txt` | `*.txt` | | | yes | Unanchored wildcard (EndsWith optimization) |
| `foo/` | `foo` | | yes | yes | Trailing `/` → directory only; no remaining `/` → NoDir |
| `foo/bar` | `foo/bar` | | | | Contains `/` → anchored to .gitignore location |
| `/foo` | `foo` | | | yes | Leading `/` stripped; no remaining `/` → NoDir, but anchored |
| `/foo/bar` | `foo/bar` | | | | Leading `/` stripped; still has `/` → anchored |
| `!foo` | `foo` | yes | | yes | Negation |
| `\!foo` | `!foo` | | | yes | Escaped `!` — literal |
| `\#comment` | `#comment` | | | yes | Escaped `#` — literal |
| `# comment` | *(skipped)* | | | | Comment line |
| *(empty)* | *(skipped)* | | | | Blank line |
| `foo   ` | `foo` | | | yes | Trailing spaces stripped |
| `foo\ ` | `foo ` | | | yes | Trailing escaped space preserved |
| `foo\\ ` | `foo\` | | | yes | `\\` = literal `\`, trailing space stripped |
| `**/foo` | `**/foo` | | | | Contains `/` → anchored |
| `foo/**` | `foo/**` | | | | Contains `/` → anchored |
| `foo/**/bar` | `foo/**/bar` | | | | Contains `/` → anchored |
| `!build/` | `build` | yes | yes | yes | Negated directory-only pattern |

**Important edge case — leading `/` + anchoring**:

The pattern `/foo` is anchored (it means "only match `foo` at this directory level"), but
after stripping the `/`, the remaining pattern `foo` has no `/`, so it gets `NoDir` flag.
This is special: despite having `NoDir`, git still knows it's anchored because the
original line had a leading `/`. In practice, the leading `/` information should be tracked
separately from the `NoDir` flag.

Actually, looking at Git's implementation more carefully: [`parse_path_pattern()`](https://github.com/git/git/blob/master/dir.c#L697-L733) strips the
`!` but does NOT strip the leading `/`. The leading `/` remains in the pattern string.
The `NODIR` flag is computed AFTER `!` is stripped but WITH the leading `/` still present.
So `/foo` → pattern is `/foo`, which contains `/`, so `NODIR` is NOT set. The leading `/`
is stripped later during matching, in [`match_pathname()`](https://github.com/git/git/blob/master/dir.c#L1352-L1415).

Corrected table:

| Input line | Stored pattern | Negative | MustBeDir | NoDir | Note |
|------------|---------------|:--------:|:---------:|:-----:|------|
| `foo` | `foo` | | | yes | No `/` → basename match |
| `*.txt` | `*.txt` | | | yes | No `/` → basename match |
| `foo/` | `foo` | | yes | yes | Trailing `/` stripped, marked MustBeDir |
| `foo/bar` | `foo/bar` | | | | Has `/` → full path match |
| `/foo` | `/foo` | | | | Has `/` → full path match; `/` stripped at match time |
| `/foo/bar` | `/foo/bar` | | | | Has `/` → full path match |
| `!foo` | `foo` | yes | | yes | Negation prefix stripped |
| `!foo/bar` | `foo/bar` | yes | | | Negation prefix stripped |
| `doc/frotz/` | `doc/frotz` | | yes | | Trailing `/` stripped; has `/` → full path |

### 2.6 Whitespace Handling Test Cases

From [`t/t0008-ignores.sh` lines 865-912](https://github.com/git/git/blob/master/t/t0008-ignores.sh#L865-L912):

| Pattern in file | Matches path | Note |
|-----------------|-------------|------|
| `trailing···` (3 spaces) | `trailing` | Trailing unescaped spaces stripped |
| `trailing\·\·` (escaped spaces) | `trailing··` | Escaped spaces preserved |
| `trailing 1 \····` (`\` then spaces) | `trailing 1 ` | `\` escapes one space, rest stripped |
| `trailing 2 \\\\` | `trailing 2 \\` | `\\` pairs are literal `\`; no trailing spaces |
| `trailing 3 \\\\·` (space after) | `trailing 3 \\` | `\\` are literal `\`; trailing space stripped |
| `trailing 4   \\\····` | `trailing 4   \ ` | Three `\`: `\\` = literal `\`, `\·` = literal space, rest stripped |
| `trailing 5 \\ \\\···` | `trailing 5 \ \ ` | `\\` = `\`, space, `\\` = `\`, `\·` = space, rest stripped |
| `trailing 6 \\a\\` | `trailing 6 \a\` | Escaped chars in middle |

(In the table above, `·` represents a space character for clarity.)

---

## Part 3: Gitignore Matching Semantics

### 3.1 Multi-Level `.gitignore` Precedence

A project can have `.gitignore` files at multiple directory levels:

```
.gitignore              ← root level
src/.gitignore          ← src/ level
src/vendor/.gitignore   ← src/vendor/ level
```

**Precedence rule**: A `.gitignore` closer to the path being checked takes priority over
one further away. Within a single `.gitignore`, later lines override earlier lines.

In practice, this means checking patterns from the deepest `.gitignore` first, then
walking up to the root `.gitignore`. The first match found (from the deepest, latest
pattern) determines the outcome.

### 3.2 Negation

A pattern starting with `!` **un-ignores** a previously ignored file:

```gitignore
*.log       # ignore all .log files
!important.log  # but keep important.log
```

**Critical limitation**: You cannot un-ignore a file if its **parent directory** is
excluded. When a directory is ignored, its contents are never even enumerated, so negation
patterns on files inside it have no effect:

```gitignore
build/          # ignore build directory
!build/output   # THIS DOES NOT WORK — build/ is already excluded
```

To work around this, you must un-ignore the directory first:

```gitignore
build/
!build/
build/*
!build/output
```

### 3.3 Directory-Only Patterns

A pattern ending with `/` only matches directories:

```gitignore
logs/    # ignores the "logs" directory (and everything inside it)
         # does NOT ignore a regular file named "logs"
```

When matching, the caller must provide information about whether the path is a directory.

### 3.4 How Patterns Match Paths

Summary of the matching flow:

```
is_ignored(path, isDir):
    for each .gitignore file, from deepest to root:
        for each pattern in file, from LAST to FIRST:
            if pattern.MustBeDir and not isDir:
                continue

            if pattern.NoDir:
                matched = wildmatch(pattern, basename(path), 0)
            else:
                rel = path relative to this .gitignore's directory
                matched = wildmatch(pattern, rel, WM_PATHNAME)

            if matched:
                if pattern.Negative:
                    return false  # un-ignored
                else:
                    return true   # ignored

    return false  # not ignored by default
```

---

## Part 4: Integration Test Cases

These test scenarios are extracted from Git's test suites and exercise the full pipeline:
parsing + matching + precedence + negation.

### 4.1 Multi-Level Setup

This is the primary test fixture from [`t/t0008-ignores.sh`](https://github.com/git/git/blob/master/t/t0008-ignores.sh#L179-L232):

```
.gitignore:
    one
    ignored-*
    top-level-dir/

a/.gitignore:
    two*
    *three

a/b/.gitignore:
    four
    five
    # this is a comment (line 3)
    six
    ignored-dir/
    # blank line follows (line 6):

    !on*
    !two
```

### 4.2 Basic Pattern Matching

| Path | isDir | Ignored? | Matched by | Note |
|------|:-----:|:--------:|-----------|------|
| `one` | no | **yes** | `.gitignore:1:one` | Exact basename match |
| `a/one` | no | **yes** | `.gitignore:1:one` | `one` has no `/` → matches basename at any depth |
| `not-ignored` | no | **no** | — | No matching pattern |
| `a/not-ignored` | no | **no** | — | No matching pattern |
| `ignored-and-untracked` | no | **yes** | `.gitignore:2:ignored-*` | Wildcard basename match |
| `a/ignored-and-untracked` | no | **yes** | `.gitignore:2:ignored-*` | Basename match at depth |
| `top-level-dir` | yes | **yes** | `.gitignore:3:top-level-dir/` | Directory-only pattern |
| `top-level-dir` | no | **no** | — | Not a directory → `top-level-dir/` doesn't match |

### 4.3 Sub-Directory Local Patterns

| Path | isDir | Ignored? | Matched by | Note |
|------|:-----:|:--------:|-----------|------|
| `a/3-three` | no | **yes** | `a/.gitignore:2:*three` | `*three` matches basename `3-three` |
| `a/three-not-this-one` | no | **no** | — | `*three` doesn't match (doesn't end with `three`) |
| `a/b/four` | no | **yes** | `a/b/.gitignore:1:four` | Exact basename match |
| `a/b/six` | no | **yes** | `a/b/.gitignore:4:six` | Line 4 (comments/blanks counted) |

### 4.4 Negation in Nested `.gitignore`

| Path | isDir | Ignored? | Matched by | Note |
|------|:-----:|:--------:|-----------|------|
| `a/b/one` | no | **no** | `a/b/.gitignore:8:!on*` | `one` matches `!on*` → un-ignored. Without this negation, would match `.gitignore:1:one` |
| `a/b/on` | no | **no** | `a/b/.gitignore:8:!on*` | `on` matches `!on*` → un-ignored |
| `a/b/two` | no | **no** | `a/b/.gitignore:9:!two` | `two` matches `!two` → un-ignored. Without this, would match `a/.gitignore:1:two*` |
| `a/b/twooo` | no | **yes** | `a/.gitignore:1:two*` | Matches `two*`; NOT negated by `!two` (not an exact match) |
| `a/b/one one` | no | **no** | `a/b/.gitignore:8:!on*` | Space in filename; `!on*` still matches |

### 4.5 Ignored Sub-Directories

| Path | isDir | Ignored? | Matched by | Note |
|------|:-----:|:--------:|-----------|------|
| `a/b/ignored-dir` | yes | **yes** | `a/b/.gitignore:5:ignored-dir/` | Directory-only pattern matches directory |
| `a/b/ignored-dir/foo` | no | **yes** | `a/b/.gitignore:5:ignored-dir/` | Parent dir is ignored → contents inherited |
| `a/b/ignored-dir/twoooo` | no | **yes** | `a/b/.gitignore:5:ignored-dir/` | Parent dir ignored, overrides `two*` |
| `a/b/ignored-dir/seven` | no | **yes** | `a/b/.gitignore:5:ignored-dir/` | Even though `a/b/ignored-dir/.gitignore` has `seven`, parent dir exclusion takes precedence |

### 4.6 Negation Cannot Override Parent Directory Exclusion

From [`t/t3001-ls-files-others-exclude.sh`](https://github.com/git/git/blob/master/t/t3001-ls-files-others-exclude.sh#L178-L194):

```gitignore
# Given these command-line exclude patterns:
--exclude="one" --exclude="!one/a.1"
```

| Path | Ignored? | Note |
|------|:--------:|------|
| `one/a.1` | **yes** | `one` directory is excluded; `!one/a.1` cannot override it |
| `one/a.2` | **yes** | Everything inside `one/` is excluded |

Conversely, negating a directory doesn't un-negate individual files:

```gitignore
--exclude="!one" --exclude="one/a.1"
```

| Path | Ignored? | Note |
|------|:--------:|------|
| `one/a.1` | **yes** | `one/a.1` pattern overrides `!one` |

### 4.7 Trailing Slash — Directory vs File

From [`t/t3001-ls-files-others-exclude.sh`](https://github.com/git/git/blob/master/t/t3001-ls-files-others-exclude.sh#L157-L170):

```gitignore
two/
```

| Path | isDir | Ignored? | Note |
|------|:-----:|:--------:|------|
| `two` (directory) | yes | **yes** | Pattern `two/` matches directory |
| `two` (regular file) | no | **no** | Pattern `two/` does NOT match regular file |

### 4.8 `**` Patterns

**Example 1**: Selective exclusion with `**` (from [t0008](https://github.com/git/git/blob/master/t/t0008-ignores.sh#L833-L848)):

```gitignore
data/**
!data/**/
!data/**/*.txt
```

| Path | isDir | Ignored? | Note |
|------|:-----:|:--------:|------|
| `data/file` | no | **yes** | Matches `data/**`; no negation matches |
| `data/data1/file1` | no | **yes** | Matches `data/**`; not a dir, not `.txt` |
| `data/data1/file1.txt` | no | **no** | Matches `data/**`, then un-ignored by `!data/**/*.txt` |
| `data/data2/file2` | no | **yes** | Matches `data/**`; not `.txt` |
| `data/data2/file2.txt` | no | **no** | Un-ignored by `!data/**/*.txt` |
| `data/data1` | yes | **no** | Un-ignored by `!data/**/` (directory negation) |

**Example 2**: `**` boundary (from [t0008](https://github.com/git/git/blob/master/t/t0008-ignores.sh#L850-L859)):

```gitignore
foo**/bar
```

| Path | isDir | Ignored? | Note |
|------|:-----:|:--------:|------|
| `foobar` | no | **no** | `foo**/bar`: `**` not at boundary → treated as `*` in pathname mode; `foobar` ≠ `foo<anything>bar` split by `/` |
| `foo/bar` | no | **yes** | Matches: `foo` then `**/bar` |

**Example 3**: `**/` matching anywhere (from [t3001](https://github.com/git/git/blob/master/t/t3001-ls-files-others-exclude.sh#L284-L293)):

```gitignore
**/a.1
```

| Path | Ignored? | Note |
|------|:--------:|------|
| `a.1` | **yes** | `**/` matches at root level |
| `one/a.1` | **yes** | `**/` matches one directory |
| `one/two/a.1` | **yes** | `**/` matches two directories |
| `three/a.1` | **yes** | `**/` matches any directory |

### 4.9 Exact Prefix Matching — `/dir/` Doesn't Match `dir-suffix/`

From [`t/t0008-ignores.sh`](https://github.com/git/git/blob/master/t/t0008-ignores.sh#L807-L831):

With `a/.gitignore` containing `/git/`:

| Path | isDir | Ignored? | Note |
|------|:-----:|:--------:|------|
| `a/git` | yes | **yes** | `/git/` matches directory `git` |
| `a/git/foo` | no | **yes** | Inside ignored directory |
| `a/git-foo` | yes | **no** | `git-foo` ≠ `git` — no prefix/substring matching |
| `a/git-foo/bar` | no | **no** | Parent not ignored |

Same behavior with `git/` (without leading `/`):

| Path | isDir | Ignored? | Note |
|------|:-----:|:--------:|------|
| `a/git` | yes | **yes** | `git/` anchored (contains `/` after stripping trailing `/`… actually `git` has no `/`, so it's unanchored basename match on directories) |
| `a/git/foo` | no | **yes** | Inside ignored directory |
| `a/git-foo` | yes | **no** | `git-foo` ≠ `git` |
| `a/git-foo/bar` | no | **no** | Parent not ignored |

### 4.10 Pattern Interaction — Complex Scenario

From [`t/t3001-ls-files-others-exclude.sh`](https://github.com/git/git/blob/master/t/t3001-ls-files-others-exclude.sh#L13-L76), a comprehensive test with multiple pattern
sources working together:

```
.gitignore:          *.1  /*.3  !*.6
one/.gitignore:      *.2  two/*.4  !*.7  *.8
one/two/.gitignore:  !*.2  !*.8
```

Plus a global exclude file: `*.7  !*.8` and a command-line exclude: `*.6`

Files in each of `.`, `one/`, `one/two/`, `three/`:
`a.1` through `a.8`

Expected untracked (not excluded) files:

| Path | Note |
|------|------|
| `a.2` | Not matched by `*.1` or `/*.3` or `*.6` |
| `a.4` | Not matched |
| `a.5` | Not matched |
| `a.8` | Matched by global `*.7` but un-ignored by `!*.8` |
| `one/a.3` | `/*.3` is anchored to root → doesn't match `one/a.3` |
| `one/a.4` | `two/*.4` only matches in `two/` subdir |
| `one/a.5` | Not matched |
| `one/a.7` | Matched by global `*.7` but un-ignored by `one/.gitignore:!*.7` |
| `one/two/a.2` | Matched by `one/.gitignore:*.2` but un-ignored by `one/two/.gitignore:!*.2` |
| `one/two/a.3` | `/*.3` anchored to root → doesn't match here |
| `one/two/a.5` | Not matched |
| `one/two/a.7` | Un-ignored by `one/.gitignore:!*.7` |
| `one/two/a.8` | Matched by `one/.gitignore:*.8` but un-ignored by `one/two/.gitignore:!*.8` |
| `three/a.2` | Not matched by root `*.1` or `/*.3` |
| `three/a.3` | `/*.3` anchored to root → doesn't match `three/a.3` |
| `three/a.4` | Not matched |
| `three/a.5` | Not matched |
| `three/a.8` | Matched by global `*.7` but un-ignored by `!*.8` |

### 4.11 Edge Cases

**UTF-8 BOM**: A `.gitignore` file starting with the UTF-8 BOM (`\xEF\xBB\xBF`) should
have the BOM stripped before parsing. The BOM must not become part of the first pattern.

**CRLF line endings**: Files with `\r\n` line endings (Windows-style) must be parsed
identically to `\n`-only files.

**Symlink `.gitignore`**: A `.gitignore` file that is a symbolic link should NOT be
followed when it's inside the working tree. This is a security measure to prevent
symlink attacks. The patterns from such a file should be silently ignored (Git logs a
warning: "unable to access .gitignore").

**Large `.gitignore` files**: Files exceeding ~100 MB should be rejected with a warning
and their patterns ignored.

---

## Part 5: Gitattributes — Differences from Gitignore

`.gitattributes` files assign attributes (like `binary`, `text`, `diff`, `merge`, etc.)
to paths using the **same pattern matching engine** as `.gitignore`. If you are implementing
both `.gitignore` and `.gitattributes`, the wildmatch engine and the core matching logic
(Parts 1-3) are fully shared. This section documents only the differences.

### 5.1 Shared Code

Both `.gitignore` and `.gitattributes` use:
- The same [`parse_path_pattern()`](https://github.com/git/git/blob/master/dir.c#L697-L733) function to parse patterns into flags
- The same [`match_basename()`](https://github.com/git/git/blob/master/dir.c#L1328-L1350) and [`match_pathname()`](https://github.com/git/git/blob/master/dir.c#L1352-L1415) functions
- The same [`wildmatch()`](https://github.com/git/git/blob/master/wildmatch.c#L286-L290) engine with the same flags

This means all wildmatch test cases (Part 1), anchoring rules, `NODIR`/`MUSTBEDIR` flag logic,
and basename-vs-fullpath matching behavior are identical.

### 5.2 Differences

| Aspect | `.gitignore` | `.gitattributes` |
|--------|-------------|-----------------|
| **Purpose** | Determines if a path is ignored (excluded from tracking) | Assigns attributes (key-value pairs) to paths |
| **Line format** | Entire line is the pattern | `pattern attr1 attr2 ...` — pattern is the first whitespace-delimited token, rest are attributes |
| **Negation (`!` prefix)** | `!pattern` un-ignores previously ignored files | **Forbidden**. Git warns "Negative patterns are ignored in git attributes" and skips the line. Use `\!` to match a literal `!`. |
| **Directory recursion** | `dir/` excludes the directory AND all its contents (Git skips traversal into excluded dirs) | `dir/` matches the directory entry only, **NOT** its contents. The docs say: "using the trailing-slash `path/` syntax is pointless in an attributes file; use `path/**` instead" |
| **Override granularity** | Last matching pattern wins — the entire match decision comes from one pattern | **Per-attribute** override — multiple patterns can match the same path, and each pattern contributes its own attributes. Later lines override earlier lines on a per-attribute basis. |
| **C-style quoting** | Not supported | Patterns starting with `"` are parsed as C-style quoted strings (supports spaces, special chars in filenames) |
| **Leading whitespace** | Part of the pattern (not stripped) | Stripped before parsing |
| **Trailing whitespace** | Stripped unless escaped with `\` (see `trim_trailing_spaces()`) | Not applicable — the pattern ends at the first whitespace; everything after is the attribute list |
| **Macro definitions** | Not applicable | Lines starting with `[attr]` define macro attributes: `[attr]binary -diff -merge -text` |
| **Max line length** | No per-line limit (only overall file size ~100MB) | 2048 bytes per line ([`ATTR_MAX_LINE_LENGTH`](https://github.com/git/git/blob/master/attr.c#L335-L338)) |
| **Attribute states** | N/A (binary: ignored or not) | `attr` (set/true), `-attr` (unset/false), `!attr` (unspecified/reset), `attr=value` (custom value) |

### 5.3 Line Parsing

A `.gitattributes` line is parsed as follows (from [`attr.c:parse_attr_line()`](https://github.com/git/git/blob/master/attr.c#L321-L410)):

1. Strip leading whitespace (spaces, tabs, `\r`, `\n`)
2. Skip blank lines and lines starting with `#` (comments)
3. If the line starts with `"`, C-style unquote the pattern
4. Otherwise, the pattern is the first whitespace-delimited token
5. If the line starts with `[attr]`, it's a macro definition (not a pattern)
6. The pattern is passed to `parse_path_pattern()` — same function as `.gitignore`
7. If `PATTERN_FLAG_NEGATIVE` is set (pattern started with `!`), reject with a warning
8. The remaining tokens are parsed as attributes

### 5.4 Matching Algorithm

```
get_attributes(path, isDir):
    for each attribute to resolve:
        value = UNSPECIFIED

    for each .gitattributes file, from deepest to root:
        for each pattern in file, from LAST to FIRST:
            if pattern.MustBeDir and not isDir:
                continue
            if pattern is a macro definition:
                continue

            if pattern.NoDir:
                matched = wildmatch(pattern, basename(path), 0)
            else:
                rel = path relative to this .gitattributes's directory
                matched = wildmatch(pattern, rel, WM_PATHNAME)

            if matched:
                for each attribute in this pattern's attribute list:
                    if attribute.value is still UNSPECIFIED:
                        attribute.value = this pattern's value
                        # (macros are expanded here too)

    return resolved attributes
```

Key difference from `.gitignore`: matching does **not** stop at the first match. All
matching patterns contribute, but only the first (= latest in file) value for each
attribute is kept.

### 5.5 Test Cases — Basic Attribute Matching

From [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L70-L158):

**Fixture:**
```
.gitattributes:
    " d "   test=d          # C-quoted pattern with spaces
     e      test=e          # leading space stripped
    f       test=f
    a/i     test=a/i        # anchored (contains /)
    onoff   test -test      # set then unset
    offon   -test test      # unset then set
    no      notest          # different attribute
    A/e/F   test=A/e/F      # mixed case anchored

a/.gitattributes:
    g       test=a/g
    b/g     test=a/b/g      # anchored within a/

a/b/.gitattributes:
    h       test=a/b/h
    d/*     test=a/b/d/*    # wildcard anchored
    d/yes   notest          # sets notest, not test
```

| Path | `test` attribute | Note |
|------|:----------------:|------|
| ` d ` | `d` | C-quoted pattern `" d "` matches filename with spaces |
| `e` | `e` | Leading space in pattern stripped |
| `f` | `f` | Basename match |
| `a/f` | `f` | `f` has no `/` → basename match at any depth |
| `a/c/f` | `f` | Works at deeper levels too |
| `a/g` | `a/g` | `a/.gitattributes:g` matches basename `g` under `a/` |
| `a/b/g` | `a/b/g` | `a/.gitattributes:b/g` (anchored) is last match |
| `b/g` | unspecified | `a/.gitattributes` patterns don't apply outside `a/` |
| `a/b/h` | `a/b/h` | `a/b/.gitattributes:h` matches |
| `a/b/d/g` | `a/b/d/*` | `d/*` in `a/b/.gitattributes` matches |
| `onoff` | unset | `test -test` — later `-test` wins (per-attr: last state) |
| `offon` | set | `-test test` — later `test` wins |
| `no` | unspecified | `no notest` sets `notest` attr, not `test` |
| `a/b/d/no` | `a/b/d/*` | `d/*` sets `test=a/b/d/*`; `notest` also set from same line |
| `a/b/d/yes` | unspecified | `d/yes notest` sets `notest` but NOT `test` |
| `a/i` | `a/i` | `a/i` has `/` → anchored to root |
| `subdir/a/i` | unspecified | `a/i` is anchored → doesn't match under `subdir/` |

### 5.6 Test Cases — Case Sensitivity

From [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L160-L204):

With `core.ignorecase=0` (case-sensitive, default):

| Path | `test` attribute | Note |
|------|:----------------:|------|
| `F` | unspecified | `f` doesn't match `F` |
| `a/F` | unspecified | Case-sensitive: no match |
| `a/b/G` | unspecified | `g` doesn't match `G` |
| `a/E/f` | `f` | `f` basename matches (directory case doesn't matter for basename patterns) |
| `A/e/F` | unspecified | Pattern `A/e/F` doesn't match because path components are checked case-sensitively |

With `core.ignorecase=1` (case-insensitive):

| Path | `test` attribute | Note |
|------|:----------------:|------|
| `F` | `f` | `f` matches `F` with casefold |
| `a/F` | `f` | Casefold basename match |
| `a/b/G` | `a/b/g` | Casefold: `G` matches `g` |
| `a/b/H` | `a/b/h` | Casefold match |
| `a/E/f` | `A/e/F` | Casefold: `A/e/F` pattern matches `a/E/f` |
| `oNoFf` | unset | Casefold: matches `onoff` pattern |

### 5.7 Test Cases — `**` Patterns

From [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L287-L322):

**`**/f` (with slashes):**
```gitattributes
**/f foo=bar
```

| Path | `foo` attribute | Note |
|------|:---------------:|------|
| `f` | `bar` | `**/` matches at root |
| `a/f` | `bar` | `**/` matches one directory |
| `a/b/f` | `bar` | `**/` matches two directories |
| `a/b/c/f` | `bar` | `**/` matches any depth |

**`a**f` (no slashes — `**` not at boundary):**
```gitattributes
a**f foo=bar
```

| Path | `foo` attribute | Note |
|------|:---------------:|------|
| `f` | unspecified | Doesn't start with `a` |
| `af` | `bar` | `a` + `**`(as `*`) + `f` |
| `axf` | `bar` | `a` + `**`(as `*`) matches `x` + `f` |
| `a/f` | unspecified | `**` not at boundary → treated as `*` in pathname mode; `*` can't cross `/` |
| `a/b/f` | unspecified | Same: `*` can't cross `/` |
| `a/b/c/f` | unspecified | Same |

### 5.8 Test Cases — Negative Pattern Rejection

From [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L276-L285):

```gitattributes
!f test=bar
```
Result: Warning "Negative patterns are ignored in git attributes". The pattern is skipped.

```gitattributes
\!f test=foo
```
Result: File named `!f` gets attribute `test=foo`. The `\!` escapes the `!`.

### 5.9 Test Cases — Prefix Matching

From [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L224-L232):

| Path | `test` attribute | Note |
|------|:----------------:|------|
| `a/g` | `a/g` | Matches `a/.gitattributes:g` |
| `a_plus/g` | unspecified | `a_plus/` is NOT confused with `a/` prefix |

### 5.10 Test Cases — Macro Expansion

From [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L489-L506):

```gitattributes
file binary
```

The built-in macro `binary` is defined as `[attr]binary -diff -merge -text`.

| Path | Attribute | Value | Note |
|------|-----------|:-----:|------|
| `file` | `binary` | set | Direct macro attribute |
| `file` | `diff` | unset | Expanded from `binary` → `-diff` |
| `file` | `merge` | unset | Expanded from `binary` → `-merge` |
| `file` | `text` | unset | Expanded from `binary` → `-text` |

### 5.11 Test Cases — Symlinks and Large Files

Same behavior as `.gitignore`:

- **Symlink `.gitattributes` in working tree**: NOT followed (warning "unable to access .gitattributes")
  — from [`t/t0003-attributes.sh:527-534`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L527-L534)
- **Large attribute line** (>2048 bytes): Warning "ignoring overly long attributes line", line skipped
  — from [`t/t0003-attributes.sh:536-543`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L536-L543)
- **Large attribute file** (>100MB): Warning "ignoring overly large gitattributes file", file skipped
  — from [`t/t0003-attributes.sh:558-563`](https://github.com/git/git/blob/master/t/t0003-attributes.sh#L558-L563)

---

## Appendix: Git Source Cross-References

### Shared (used by both `.gitignore` and `.gitattributes`)

| Concept | Git source location |
|---------|-------------------|
| Wildmatch engine | [`wildmatch.c:59-290`](https://github.com/git/git/blob/master/wildmatch.c#L59-L290) (function `dowild`) |
| Wildmatch public API | [`wildmatch.h`](https://github.com/git/git/blob/master/wildmatch.h) |
| Glob special chars (`*?[\`) | [`sane-ctype.h:47`](https://github.com/git/git/blob/master/sane-ctype.h#L47) (`is_glob_special`) |
| Pattern parsing | [`dir.c:697-733`](https://github.com/git/git/blob/master/dir.c#L697-L733) (`parse_path_pattern`) |
| Basename matching | [`dir.c:1328-1350`](https://github.com/git/git/blob/master/dir.c#L1328-L1350) (`match_basename`) |
| Full path matching | [`dir.c:1352-1415`](https://github.com/git/git/blob/master/dir.c#L1352-L1415) (`match_pathname`) |
| Pattern flags | [`dir.h:52-55`](https://github.com/git/git/blob/master/dir.h#L52-L55) (`PATTERN_FLAG_*`) |
| Wildmatch test suite | [`t/t3070-wildmatch.sh`](https://github.com/git/git/blob/master/t/t3070-wildmatch.sh) |

### Gitignore-specific

| Concept | Git source location |
|---------|-------------------|
| Trailing space stripping | [`dir.c:1029-1050`](https://github.com/git/git/blob/master/dir.c#L1029-L1050) (`trim_trailing_spaces`) |
| File buffer parsing (BOM, CRLF, comments) | [`dir.c:1226-1254`](https://github.com/git/git/blob/master/dir.c#L1226-L1254) (`add_patterns_from_buffer`) |
| Pattern list scanning (last match wins) | [`dir.c:1423-1469`](https://github.com/git/git/blob/master/dir.c#L1423-L1469) (`last_matching_pattern_from_list`) |
| Multi-source precedence | [`dir.c:1624-1643`](https://github.com/git/git/blob/master/dir.c#L1624-L1643) (`last_matching_pattern_from_lists`) |
| Public exclude check API | [`dir.c:1831-1839`](https://github.com/git/git/blob/master/dir.c#L1831-L1839) (`is_excluded`) |
| Gitignore test suite | [`t/t0008-ignores.sh`](https://github.com/git/git/blob/master/t/t0008-ignores.sh) |
| ls-files exclude tests | [`t/t3001-ls-files-others-exclude.sh`](https://github.com/git/git/blob/master/t/t3001-ls-files-others-exclude.sh) |
| Official gitignore docs | [`Documentation/gitignore.adoc`](https://github.com/git/git/blob/master/Documentation/gitignore.adoc) |

### Gitattributes-specific

| Concept | Git source location |
|---------|-------------------|
| Attribute line parsing | [`attr.c:321-410`](https://github.com/git/git/blob/master/attr.c#L321-L410) (`parse_attr_line`) |
| Calls shared `parse_path_pattern()` | [`attr.c:385-388`](https://github.com/git/git/blob/master/attr.c#L385-L388) |
| Negative pattern rejection | [`attr.c:389-393`](https://github.com/git/git/blob/master/attr.c#L389-L393) |
| Attribute value parsing | [`attr.c:276-319`](https://github.com/git/git/blob/master/attr.c#L276-L319) (`parse_attr`) |
| Path matching (calls shared funcs) | [`attr.c:1045-1066`](https://github.com/git/git/blob/master/attr.c#L1045-L1066) (`path_matches`) |
| Per-attribute fill logic | [`attr.c:1092-1115`](https://github.com/git/git/blob/master/attr.c#L1092-L1115) (`fill_one`) |
| Stack scanning | [`attr.c:1117-1135`](https://github.com/git/git/blob/master/attr.c#L1117-L1135) (`fill`) |
| Macro expansion | [`attr.c:1109-1110`](https://github.com/git/git/blob/master/attr.c#L1109-L1110) (in `fill_one`) |
| Built-in macros | [`attr.c:638-641`](https://github.com/git/git/blob/master/attr.c#L638-L641) (`binary`) |
| Data structures | [`attr.h:253-281`](https://github.com/git/git/blob/master/attr.h#L253-L281) (`struct pattern`, `struct match_attr`) |
| Gitattributes test suite | [`t/t0003-attributes.sh`](https://github.com/git/git/blob/master/t/t0003-attributes.sh) |
| Archive attr pattern tests | [`t/t5002-archive-attr-pattern.sh`](https://github.com/git/git/blob/master/t/t5002-archive-attr-pattern.sh) |
| Official gitattributes docs | [`Documentation/gitattributes.adoc`](https://github.com/git/git/blob/master/Documentation/gitattributes.adoc) |
