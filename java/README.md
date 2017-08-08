# enry-java

### Requirements

* `sbt`
* `Java` (tested with Java 1.8)
* `maven` install and on the PATH (only for local usage)
* `Go` (only for building the shared objects for your operating system)

### Generate jar with Java bindings and shared libraries

You need to do this before exporting the jar and/or testing.

```
make
```

This will download JNAerator and build its jar to generate the code from the `libenry.h` header file (hence the need for `mvn` installed), it will be placed under `lib`.
The shared libraries for your operating system will be built if needed and copied inside the `shared` directory.

For IntelliJ and other IDEs remember to mark `shared` folder as sources and add `lib/enry.jar` as library. If you use `sbt` from the command line directly that's already taken care of.

### Run tests

```
make test
```

### Export jar

```
make package
```

Jar will be located in `./target/enry-java-assembly-X.X.X.jar`.
