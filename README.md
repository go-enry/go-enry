# enry [![GoDoc](https://godoc.org/gopkg.in/src-d/enry.v1?status.svg)](https://godoc.org/gopkg.in/src-d/enry.v1) [![Build Status](https://travis-ci.org/src-d/enry.svg?branch=master)](https://travis-ci.org/src-d/enry) [![codecov](https://codecov.io/gh/src-d/enry/branch/master/graph/badge.svg)](https://codecov.io/gh/src-d/enry)

File programming language detector and toolbox to ignore binary or vendored files. *enry*, started as a port to _Go_ of the original [linguist](https://github.com/github/linguist) _Ruby_ library, that has an improved *performance of 100x*.


Installation
------------

The recommended way to install simple-linguist

```
go get gopkg.in/src-d/enry.v1/...
```

To build enry's CLI you must run

    make build-cli

it generates a binary in the project's root directory called `enry`. You can move this binary to anywhere in your `PATH`.

Examples
--------

```go
lang, safe := enry.GetLanguageByExtension("foo.go")
fmt.Println(lang)
// result: Go

lang, safe := enry.GetLanguageByContent("foo.m", "<matlab-code>")
fmt.Println(lang)
// result: Matlab

lang, safe := enry.GetLanguageByContent("bar.m", "<objective-c-code>")
fmt.Println(lang)
// result: Objective-C

// all strategies together
lang := enry.GetLanguage("foo.cpp", "<cpp-code>")
```

Note the returned boolean value "safe" is set either to true, if there is only one possible language detected or, to false otherwise.

To get a list of possible languages for a given file, you can use the plural version of the detecting functions.

```go
langs := enry.GetLanguages("foo.h",  "<cpp-code>")
// result: []string{"C++", "C"}

langs := enry.GetLanguagesByExtension("foo.asc", "<content>", nil)
// result: []string{"AGS Script", "AsciiDoc", "Public Key"}

langs := enry.GetLanguagesByFilename("Gemfile", "<content>", []string{})
// result: []string{"Ruby"}
```


CLI
-----------------

You can use enry as a command,

```bash
$ enry --help
enry, A simple (and faster) implementation of github/linguist
usage: enry <path>
              enry <path> [--json] [--breakdown]
              enry [--json] [--breakdown]
```

and it will return an output similar to *linguist*'s output,

```bash
$ enry
11.11%    Gnuplot
22.22%    Ruby
55.56%    Shell
11.11%    Go
```

but not only the output, also its flags are the same as *linguist*'s ones,

```bash
$ enry --breakdown
11.11%    Gnuplot
22.22%    Ruby
55.56%    Shell
11.11%    Go

Gnuplot
plot-histogram.gp

Ruby
linguist-samples.rb
linguist-total.rb

Shell
parse.sh
plot-histogram.sh
run-benchmark.sh
run-slow-benchmark.sh
run.sh

Go
parser/main.go
```

even the JSON flag,

```bash
$ enry --json
{"Gnuplot":["plot-histogram.gp"],"Go":["parser/main.go"],"Ruby":["linguist-samples.rb","linguist-total.rb"],"Shell":["parse.sh","plot-histogram.sh","run-benchmark.sh","run-slow-benchmark.sh","run.sh"]}
```

Note that even if enry's CLI is compatible with linguist's, its main point is that, contrary to linguist, **_enry doesn't need a git repository to work!_**


Development
-----------

*enry* re-uses parts of original [linguist](https://github.com/github/linguist) especially data in `languages.yml` to generate internal data structures. In oreder to update to latest upstream run

    make clean code-generate

To run the tests

    make test


Why Enry?
---------

In the movie [My Fair Lady](https://en.wikipedia.org/wiki/My_Fair_Lady), [Professor Henry Higgins](http://www.imdb.com/character/ch0011719/?ref_=tt_cl_t2) is one of the main characters. Henry is a linguist and at the very beginning of the movie enjoys guessing the nationality of people based on their accent.

`Enry Iggins` is how [Eliza Doolittle](http://www.imdb.com/character/ch0011720/?ref_=tt_cl_t1), [pronounces](https://www.youtube.com/watch?v=pwNKyTktDIE) the name of the Professor during the first half of the movie.


License
-------

MIT, see [LICENSE](LICENSE)
