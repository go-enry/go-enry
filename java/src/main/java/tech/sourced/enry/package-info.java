/**
 * Java bindings to the <a href="https://github.com/go-enry/go-enry">go-enry</a> parsing library.
 *
 * <h2 id="requirements">Requirements</h2>
 *
 * <ul>
 * <li>
 *     JDK 22 (for the
 *     <a href="https://docs.oracle.com/en/java/javase/22/core/foreign-function-and-memory-api.html">
 *         Foreign Function and Memory API
 *     </a>)
 * </li>
 * <li>Go-enry shared library</li>
 * <li>Generated bindings for languages</li>
 * <li>Shared libraries for languages</li>
 * </ul>
 *
 * <h2 id="usage">Basic Usage</h2>
 *
 * {@snippet lang = java:
 * Enry.getLanguage("Main.java", "".getBytes());
 *}
 *
 * <h2 id="libraries">Library Loading</h2>
 *
 * There are three ways to load the shared libraries:
 *
 * <ol>
 * <li>
 *     The libraries can be installed in the OS-specific library search path or in
 *     {@systemProperty java.library.path}. The search path can be amended using the
 *     {@code LD_LIBRARY_PATH} environment variable on Linux, {@code DYLD_LIBRARY_PATH}
 *     on macOS, or {@code PATH} on Windows. The libraries will be loaded automatically by
 *     {@link java.lang.foreign.SymbolLookup#libraryLookup(String, java.lang.foreign.Arena)
 *     SymbolLookup.libraryLookup(String, Arena)}.
 * </li>
 * <li>
 *     The libraries can be loaded manually by calling
 *     {@link java.lang.System#loadLibrary(String) System.loadLibrary(String)},
 *     if the library is installed in {@systemProperty java.library.path},
 *     or {@link java.lang.System#load(String) System.load(String)}.
 * </li>
 * <li>
 *     The libraries can be loaded manually by registering a custom implementation of
 *     {@link tech.sourced.enry.NativeLibraryLookup NativeLibraryLookup}.
 *     This can be used, for example, to load libraries from inside a JAR file.
 * </li>
 * </ol>
 */
package tech.sourced.enry;
