# enry [![GoDoc](https://godoc.org/gopkg.in/src-d/enry.v1?status.svg)](https://godoc.org/gopkg.in/src-d/enry.v1) [![Build Status](https://travis-ci.org/src-d/enry.svg?branch=master)](https://travis-ci.org/src-d/simple-linguist)

File programming language detector and toolbox to ignore binary or vendored files. *enry*, started as a port to _Go_ of the original [linguist](https://github.com/github/linguist) _Ruby_ library, that has an improved *performance of 100x*.


Installation
------------

The recommended way to install simple-linguist

```
go get gopkg.in/src-d/enry.v1/...
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

Why Enry?
---------
In the movie [My Fair Lady](https://en.wikipedia.org/wiki/My_Fair_Lady), [Professor Henry Higgins](http://www.imdb.com/character/ch0011719/?ref_=tt_cl_t2) is one of the main characters. Henry is a linguist and at the very beginning of the movie enjoys guessing the nationality of people based on their accent.

`Enry Iggins` is how [Eliza Doolittle](http://www.imdb.com/character/ch0011720/?ref_=tt_cl_t1), [pronounces](https://www.youtube.com/watch?v=pwNKyTktDIE) the name of the Professor during the first half of the movie.


License
-------

MIT, see [LICENSE](LICENSE)
