name := "enry-java"
organization := "tech.sourced"
version := "1.0"

crossPaths := false
autoScalaLibrary := false
publishMavenStyle := true
exportJars := true

libraryDependencies += "com.novocode" % "junit-interface" % "0.11" % Test

unmanagedBase := baseDirectory.value / "lib"
unmanagedClasspath in Test += baseDirectory.value / "shared"
unmanagedClasspath in Runtime += baseDirectory.value / "shared"
unmanagedClasspath in Compile += baseDirectory.value / "shared"
testOptions += Tests.Argument(TestFrameworks.JUnit)

