# enry-java

## Usage

`enry-java` package is available through [maven central],
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

* `maven`
* `Java` (tested with Java 22)
* `Jextract` mechanically generates Java bindings from native library headers
* `Go` (only for building the shared objects for your operating system)

### Generate jar with Java bindings and shared libraries

You need to do this before exporting the jar and/or testing.

```bash
make
```

The shared libraries for your operating system will be built if needed and copied inside the `shared` directory.

### Run tests

```bash
make test
```

### Export jar

```bash
make package
```

Will build fatJar under `./target/enry-java-X.X.X-SNAPSHOT.jar`.

[maven central]: http://search.maven.org/#search%7Cga%7C1%7Ca%3A%22enry-java%22
