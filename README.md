# enry [![GoDoc](https://godoc.org/gopkg.in/src-d/enry.v1?status.svg)](https://godoc.org/gopkg.in/src-d/enry.v1) [![Build Status](https://travis-ci.org/src-d/enry.svg?branch=master)](https://travis-ci.org/src-d/enry) [![codecov](https://codecov.io/gh/src-d/enry/branch/master/graph/badge.svg)](https://codecov.io/gh/src-d/enry)

File programming language detector and toolbox to ignore binary or vendored files. *enry*, started as a port to _Go_ of the original [linguist](https://github.com/github/linguist) _Ruby_ library, that has an improved *2x performance*.


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
------------

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
------------

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
------------

*enry* re-uses parts of original [linguist](https://github.com/github/linguist) especially data in `languages.yml` to generate internal data structures. In oreder to update to latest upstream run

    make clean code-generate

To run the tests

    make test


Divergences from linguist
------------

Using [linguist/samples](https://github.com/github/linguist/tree/master/samples) as a set against run tests the following issues were found:
* with [hello.ms](https://github.com/github/linguist/blob/master/samples/Unix%20Assembly/hello.ms) we can't detect the language (Unix Assembly) because we don't have a matcher in contentMatchers (content.go) for Unix Assembly. Linguist uses this [regexp](https://github.com/github/linguist/blob/master/lib/linguist/heuristics.rb#L300) in its code,

    `elsif /(?<!\S)\.(include|globa?l)\s/.match(data) || /(?<!\/\*)(\A|\n)\s*\.[A-Za-z][_A-Za-z0-9]*:/.match(data.gsub(/"([^\\"]|\\.)*"|'([^\\']|\\.)*'|\\\s*(?:--.*)?\n/, ""))`

    which we can't port.

* all files for SQL language fall to the classifier because we don't parse this [disambiguator expresion](https://github.com/github/linguist/blob/master/lib/linguist/heuristics.rb#L433) for `*.sql` files right. This expression doesn't comply with the pattern for the rest of [heuristics.rb](https://github.com/github/linguist/blob/master/lib/linguist/heuristics.rb) file.


Benchmarks
------------

Enry's language detection has been compared with Linguist's language detection. In order to do that, linguist's project directory [*linguist/samples*](https://github.com/github/linguist/tree/master/samples) was used as a set of files to run benchmarks against.

 Following results were obtained:

![histogram](https://raw.githubusercontent.com/src-c/enry/master/benchmarks/histogram/distribution.jpg)

The histogram represents the number of files for which spent time in language detection was in the range of the time interval indicated in x axis.

So reviewing the comparison enry/linguist, you can see the most of the files were detected in less time than linguist does.

We detected some few cases enry turns slower than linguist. This is due to Golang's regexp engine being slower than Ruby's, which uses [oniguruma](https://github.com/kkos/oniguruma) library, written in C.

You can find scripts and additional information (as software and hardware used, and benchmarks' results per sample file) in [*benchmarks*](https://github.com/src-d/enry/tree/master/benchmarks) directory.

If you want to reproduce the same experiment you can run:

    benchmarks/run.sh

from the root's project directory and It runs benchmarks for enry and linguist, parse the output, create csv files and create a histogram (you must have installed [gnuplot](http://gnuplot.info) in your system to get the histogram). It can take to much time, so to run local benchmarks to take a quick look you can run either:

    make benchmarks

to get time averages for main detection function and strategies for the whole samples set or:

    make benchmarks-samples

if you want see measures by sample file


Why Enry?
------------

In the movie [My Fair Lady](https://en.wikipedia.org/wiki/My_Fair_Lady), [Professor Henry Higgins](http://www.imdb.com/character/ch0011719/?ref_=tt_cl_t2) is one of the main characters. Henry is a linguist and at the very beginning of the movie enjoys guessing the nationality of people based on their accent.

`Enry Iggins` is how [Eliza Doolittle](http://www.imdb.com/character/ch0011720/?ref_=tt_cl_t1), [pronounces](https://www.youtube.com/watch?v=pwNKyTktDIE) the name of the Professor during the first half of the movie.


License
------------

Apache License, Version 2.0. See [LICENSE](LICENSE)
