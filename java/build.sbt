import scala.sys.process._

name := "enry-java"
organization := "tech.sourced"
enablePlugins(GitVersioning)
git.useGitDescribe := true
git.gitDescribePatterns := Seq("v[0-9]*")

homepage := Some(url("https://github.com/go-enry/go-enry"))
licenses += ("Apache-2.0", url("https://www.apache.org/licenses/LICENSE-2.0"))
crossPaths := false
autoScalaLibrary := false
publishMavenStyle := true

libraryDependencies += "net.java.dev.jna" % "jna" % "5.19.1"
libraryDependencies += "com.github.sbt" % "junit-interface" % "0.13.3" % Test
javacOptions ++= Seq("--release", "8", "-encoding", "UTF-8")
Test / fork := true
Test / javaOptions += "-Djna.nosys=true"

// JNA extracts the library for the running platform from the jar resource.
// Always rebuild: Go's own cache tracks detector/data changes outside java/.
Compile / resourceGenerators += Def.task {
  val root = baseDirectory.value.getParentFile
  val go = sys.env.getOrElse("GO", "go")
  val os = Process(Seq(go, "env", "GOOS"), root).!!.trim
  val arch = Process(Seq(go, "env", "GOARCH"), root).!!.trim
  val jnaArch = arch match {
    case "amd64" => "x86-64"
    case "arm64" => "aarch64"
    case other => sys.error(s"Unsupported native architecture: $other")
  }
  val suffix = os match {
    case "darwin" => "dylib"
    case "linux" => "so"
    case other => sys.error(s"Unsupported native operating system: $other")
  }
  val directory = (Compile / resourceManaged).value / s"$os-$jnaArch"
  IO.createDirectory(directory)
  val library = directory / s"libenry.$suffix"
  val status = Process(Seq(go, "build", "-buildmode=c-shared", "-o", library.getAbsolutePath,
    "./shared"), root, "CGO_ENABLED" -> "1").!
  if (status != 0) sys.error("Failed to build the enry native library")
  Seq(library)
}.taskValue

// The assembly contains JNA and this platform's native library.
assembly / assemblyMergeStrategy := {
  case PathList("META-INF", "versions", _, "module-info.class") => MergeStrategy.discard
  case path => (assembly / assemblyMergeStrategy).value(path)
}
