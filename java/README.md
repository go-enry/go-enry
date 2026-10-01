# enry-java

## Usage

`enry-java` package is available through [maven central](http://search.maven.org/#search%7Cga%7C1%7Ca%3A%22enry-java%22),
so it be used easily added as a dependency in various package management systems.
Examples of how to handle it for most commons systems are included below,
for other systems just look at maven central's dependency information.

### Apache Maven

```xml
<dependency>
    <groupId>tech.sourced</groupId>
    <artifactId>enry-java</artifactId>
    <version>${enry_version}</version>
</dependency>
```

### Scala SBT

```scala
libraryDependencies += "tech.sourced" % "enry-java" % enryVersion
```

## Build

### Requirements

* JDK 17 or later to build (the Java API targets Java 8 bytecode)
* Go 1.26.x and a C compiler
* `curl` and `shasum` for the checksum-verified sbt launcher
* Linux or macOS on x86_64 or arm64

From `java/`:

```bash
make test
make package
```

The build compiles `../shared` using the current C ABI and packages its native
library under JNA's platform-specific resource directory. JNA extracts the
matching library from the jar at runtime; no external `libenry` installation or
JNAerator-generated jar is needed. `make package` creates
`target/enry-java-assembly-*.jar`, including JNA and the current platform's native
library. Build separately on each target platform; one build is not a universal
jar. `./sbt publishLocal` installs the platform build locally. Maven Central
publication is not automated by this build.

All string arguments use UTF-8. Null strings and byte arrays are treated as
empty. NUL bytes are supported in content arrays, but rejected in string
arguments because the C ABI uses NUL-terminated strings. Native string and array
results are copied and freed, and `Guess.safe` retains the detector's ambiguity
information through the `WithSafety` exports.

The Maven coordinates above describe the existing published package; this
binding repair must be released separately before those artifacts include it.
