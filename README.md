# simple-linguist

File language detector and toolbox to ignore binary or vendored files. *simple-linguist*, is our port to _Go_ of the original [lignuist](https://github.com/github/linguist) _Ruby_ library, with fewer precision in arcane languages but with an improved *performance of 100x*.


Installation
------------

The recommended way to install simple-linguist

```
go get gopkg.in/src-d/simple-linguist.v1/...
```


Examples
--------

```go
lang, _ := GetLanguageByExtension("foo.go")
fmt.Println(lang)
// result: Go

lang, _ = GetLanguageByContent("foo.m", "<matlab-code>")
fmt.Println(lang)
// result: Matlab

lang, _ = GetLanguageByContent("bar.m", "<pbjective-c-code>")
fmt.Println(lang)
// result: Objective-C
```