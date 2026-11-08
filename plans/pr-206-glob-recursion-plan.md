# PR #206 plan: bound `.gitattributes` glob backtracking

## Problem

PR #206 implements `.gitattributes` pattern matching in `gitattributes.go`.

The current `globMatchInternal(pattern, str, prev)` implementation recursively explores wildcard splits. Patterns with many `*` segments can cause exponential backtracking on non-matching input.

The maintainer-provided example is:

```gitattributes
*x*x*x*x*x*x*x*x*x*xz linguist-vendored
```

matched against a long basename such as:

```text
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.go
```

This can take seconds of CPU time and should be bounded before merging.

## Goals

- Preserve current matching semantics.
- Bound runtime for malicious wildcard-heavy patterns.
- Keep the implementation local to `gitattributes.go`.
- Avoid new public API.
- Add tests that would have caught the PR review example.

## Non-goals

- Do not replace the matcher with go-git.
- Do not implement nested `.gitattributes`.
- Do not implement case-folding.
- Do not change documented matching behavior except for pathological runtime safety.

## Current code path

File: `gitattributes.go`

- `matchGitPattern(pattern, path string) bool`
- `globMatch(pattern, str string) bool`
- `globMatchInternal(pattern, str string, prev byte) bool`

The recursive calls happen in:

- `**` path-boundary branch, looping over candidate slash boundaries.
- `*` branch, looping over every candidate split before `/`.

Both branches can revisit the same `(pattern suffix, string suffix, prev)` states many times.

## Recommended option: memoized recursive matcher

Use dynamic programming memoization keyed by pattern index, string index, and `prev`. This keeps the existing recursive structure readable while preventing repeated work.

### Implementation shape

File: `gitattributes.go`

Estimated production diff: 60-100 LoC.

1. Replace the string-slicing recursive helper with an index-based helper.

Keep the public/private entry function signature:

```go
func globMatch(pattern, str string) bool
```

Change it to:

```go
func globMatch(pattern, str string) bool {
    memo := make(map[globMatchState]bool)
    seen := make(map[globMatchState]bool)
    return globMatchInternal(pattern, str, 0, 0, 0, memo, seen)
}
```

2. Add a private state type:

```go
type globMatchState struct {
    patternIndex int
    stringIndex  int
    prev         byte
}
```

3. Change `globMatchInternal` to use indices:

```go
func globMatchInternal(pattern, str string, pi, si int, prev byte, memo map[globMatchState]bool, seen map[globMatchState]bool) bool {
    state := globMatchState{patternIndex: pi, stringIndex: si, prev: prev}
    if seen[state] {
        return memo[state]
    }
    seen[state] = true

    matched := globMatchInternalUncached(pattern, str, pi, si, prev, memo, seen)
    memo[state] = matched
    return matched
}
```

`globMatchInternalUncached` can contain the existing switch logic, translated from slicing to index movement.

4. Translation rules:

Current:

```go
for len(pattern) > 0 {
    switch pattern[0] {
```

New:

```go
for pi < len(pattern) {
    switch pattern[pi] {
```

Current:

```go
pattern = pattern[1:]
str = str[1:]
```

New:

```go
pi++
si++
```

Current recursive call:

```go
globMatchInternal(rest, str[i:], restPrev)
```

New recursive call:

```go
globMatchInternal(pattern, str, restIndex, si+i, restPrev, memo, seen)
```

Current end condition:

```go
return len(str) == 0
```

New:

```go
return si == len(str)
```

5. Update bracket handling.

The existing helper `findClosingBracket(pattern string) int` returns an index relative to the passed substring. To reduce churn, either:

- Keep using it on `pattern[pi:]` and convert the returned offset with `end := pi + relEnd`.
- Or add `findClosingBracketFrom(pattern string, pi int) int`.

Lowest-risk patch: use the first approach.

6. Add an optional internal state budget only if reviewers want a literal "limit".

The memoized algorithm already bounds states to approximately:

```text
len(pattern) * len(str) * distinct(prev)
```

For a defense-in-depth cap:

```go
const maxGlobMatchStates = 100000
```

If `len(seen) > maxGlobMatchStates`, return `false`.

This is conservative but can create false negatives for very large legitimate patterns. Prefer memoization alone unless review specifically asks for a hard cap.

### Why this is preferable to a hard recursion-depth limit

A recursion-depth limit stops stack growth, but it does not directly stop repeated work before the limit is reached. It can also reject valid deep patterns depending on where the limit is set.

Memoization solves the actual bug: repeated exploration of the same wildcard state.

## Test changes

File: `gitattributes_edge_test.go`

Estimated test diff: 35-70 LoC.

### Add the maintainer-provided pathological case

Use a bounded-time test without making the suite flaky.

```go
func TestEdgeGlobMatchPathologicalBacktracking(t *testing.T) {
    pattern := "*x*x*x*x*x*x*x*x*x*xz"
    path := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.go"

    start := time.Now()
    got := matchGitPattern(pattern, path)
    elapsed := time.Since(start)

    assert.False(t, got)
    assert.Less(t, elapsed, 100*time.Millisecond)
}
```

The threshold should be generous enough for CI. If maintainers dislike timing assertions, use a benchmark-style test with a goroutine timeout:

```go
done := make(chan bool, 1)
go func() {
    done <- matchGitPattern(pattern, path)
}()

select {
case got := <-done:
    assert.False(t, got)
case <-time.After(500 * time.Millisecond):
    t.Fatal("pathological glob match exceeded timeout")
}
```

The direct elapsed assertion is simpler but can be more sensitive to slow builders.

### Add no-regression tests for normal wildcard semantics

Keep these near existing `globMatch` tests:

```go
{"*.go", "main.go", true}
{"*x*z", "xxz", true}
{"*x*z", "xx.go", false}
{"foo**bar", "foobazbar", true}
{"foo**bar", "foo/baz/bar", false}
{"**/*.go", "src/main.go", true}
{"**/*.go", "main.go", true}
```

### Add the unresolved inline review nit

In `TestEdgeTrailingSlashDirectory`, add:

```go
{"vendor/", "src/vendor/", true, "nested directory matches non-anchored pattern"}
```

This addresses the unresolved review thread on `gitattributes_edge_test.go`.

## Validation

Run:

```bash
go test ./...
go test -run 'TestEdgeGlobMatchPathologicalBacktracking|TestEdgeTrailingSlashDirectory|TestMatchGitPattern|TestEdgeDoubleStar' -count=20
```

Optional micro-benchmark for reviewer confidence:

```go
func BenchmarkGlobMatchPathological(b *testing.B) {
    pattern := "*x*x*x*x*x*x*x*x*x*xz"
    path := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.go"
    for i := 0; i < b.N; i++ {
        _ = matchGitPattern(pattern, path)
    }
}
```

Run:

```bash
go test -run '^$' -bench 'BenchmarkGlobMatchPathological' -benchmem
```

Expected result: the benchmark should complete quickly and should not grow exponentially as the input length increases.

## Option B: recursion and step budget without memoization

Estimated production diff: 35-60 LoC.

Add a mutable budget object passed through recursion:

```go
type globMatchBudget struct {
    steps int
}
```

Each recursive call increments `steps`. If it exceeds a fixed threshold, return false.

### Pros

- Smaller rewrite than index-based memoization.
- Directly matches the phrase "recursion limit".

### Cons

- Does not remove repeated work.
- May still burn CPU until the budget is exhausted.
- Threshold selection is arbitrary.
- Can create false negatives for legitimate large patterns.

Recommendation: use only as defense-in-depth on top of memoization, not as the primary fix.

## Option C: iterative glob matcher

Estimated production diff: 120-200 LoC.

Rewrite the matcher as an explicit NFA-style state machine:

- States contain pattern index, string index, and previous pattern byte.
- Use a queue or stack.
- Track visited states.
- Return true when an accepting state is reached.

### Pros

- No recursion.
- Natural visited-state bounding.
- Very explicit control over memory and work.

### Cons

- Larger rewrite.
- More review burden.
- Higher risk of introducing semantic regressions in a PR that is otherwise close to merge.

Recommendation: defer unless maintainers want a broader matcher rewrite.

## Reviewer-facing summary

The smallest robust fix is memoization. It preserves the current matcher behavior while bounding malicious wildcard-heavy input. A hard recursion limit alone is less correct because it limits stack depth, not repeated state exploration.

Rough total size for the recommended patch:

- `gitattributes.go`: 60-100 LoC changed or added.
- `gitattributes_edge_test.go`: 35-70 LoC added.
- Total: about 95-170 LoC, depending on whether the helper is factored as a separate uncached function and whether a hard state cap is included.
