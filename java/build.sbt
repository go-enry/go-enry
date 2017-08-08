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

lazy val buildNative = taskKey[Unit]("builds native code")

buildNative := {
  val res = "make"!;
  if (res != 0) throw new RuntimeException("unable to generate shared libraries and native jar bindings")
}

test := {
  buildNative.value
  (test in Test).value
}

compile := {
  buildNative.value
  (compile in Compile).value
}

assembly := {
  buildNative.value
  assembly.value
}
