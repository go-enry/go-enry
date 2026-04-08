package tech.sourced.enry;

import tech.sourced.enry.internal.*;

import java.lang.foreign.Arena;
import java.lang.foreign.MemorySegment;

import static tech.sourced.enry.GoUtils.*;

public class Enry {
    public static final Guess unknownLanguage = new Guess("", false);

    /**
     * Returns the language of the given file based on the filename and its
     * contents.
     *
     * @param filename name of the file with the extension
     * @param content  array of bytes with the contents of the file (the code)
     * @return the guessed language
     */
    public static synchronized String getLanguage(String filename, byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaString(GoEnry.GetLanguage(
                arena,
                toGoString(arena, filename),
                toGoByteSlice(arena, content)
            ));
        }
    }

    /**
     * Returns detected language by its content.
     * If there are more than one possible language, it returns the first
     * language in alphabetical order and safe to false.
     *
     * @param filename name of the file with the extension
     * @param content  of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByContent(String filename, byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByContent(
                arena,
                toGoString(arena, filename),
                toGoByteSlice(arena, content)
            );
            return new Guess(
                toJavaString(GetLanguageByContent_return.r0(res)),
                toJavaBool(GetLanguageByContent_return.r1(res))
            );
        }
    }

    /**
     * Returns detected language by its emacs modeline.
     * If there are more than one possible language, it returns the first
     * language in alphabetical order and safe to false.
     *
     * @param content of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByEmacsModeline(byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByEmacsModeline(arena, toGoByteSlice(arena, content));
            return new Guess(
                toJavaString(GetLanguageByEmacsModeline_return.r0(res)),
                toJavaBool(GetLanguageByEmacsModeline_return.r1(res))
            );
        }
    }

    /**
     * Returns detected language by the extension of the filename.
     * If there are more than one possible languages, it returns
     * the first language in alphabetical order and safe to false.
     *
     * @param filename of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByExtension(String filename) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByExtension(arena, toGoString(arena, filename));
            return new Guess(
                toJavaString(GetLanguageByExtension_return.r0(res)),
                toJavaBool(GetLanguageByExtension_return.r1(res))
            );
        }
    }

    /**
     * Returns detected language by its shebang.
     * If there are more than one possible language, it returns the first
     * language in alphabetical order and safe to false.
     *
     * @param content of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByShebang(byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByShebang(arena, toGoByteSlice(arena, content));
            return new Guess(
                toJavaString(GetLanguageByShebang_return.r0(res)),
                toJavaBool(GetLanguageByShebang_return.r1(res))
            );
        }
    }

    /**
     * Returns detected language by its filename.
     * If there are more than one possible language, it returns the first
     * language in alphabetical order and safe to false.
     *
     * @param filename of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByFilename(String filename) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByFilename(arena, toGoString(arena, filename));
            return new Guess(
                toJavaString(GetLanguageByFilename_return.r0(res)),
                toJavaBool(GetLanguageByFilename_return.r1(res))
            );
        }
    }

    /**
     * Returns detected language by its modeline.
     * If there are more than one possible language, it returns the first
     * language in alphabetical order and safe to false.
     *
     * @param content of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByModeline(byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByModeline(arena, toGoByteSlice(arena, content));
            return new Guess(
                toJavaString(GetLanguageByModeline_return.r0(res)),
                toJavaBool(GetLanguageByModeline_return.r1(res))
            );
        }
    }

    /**
     * Returns detected language by its vim modeline.
     * If there are more than one possible language, it returns the first
     * language in alphabetical order and safe to false.
     *
     * @param content of the file
     * @return guessed result
     */
    public static synchronized Guess getLanguageByVimModeline(byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment res = GoEnry.GetLanguageByVimModeline(arena, toGoByteSlice(arena, content));
            return new Guess(
                toJavaString(GetLanguageByVimModeline_return.r0(res)),
                toJavaBool(GetLanguageByVimModeline_return.r1(res))
            );
        }
    }

    /**
     * Returns all the possible extensions for a file in the given language.
     *
     * @param language to get extensions from
     * @return extensions
     */
    public static synchronized String[] getLanguageExtensions(String language) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment result = GoSlice.allocate(arena);
            GoEnry.GetLanguageExtensions(toGoString(arena, language), result);
            return toJavaStringArray(result);
        }
    }

    /**
     * Returns all possible languages for the given file.
     *
     * @param filename of the file
     * @param content  of the file
     * @return all possible languages
     */
    public static synchronized String[] getLanguages(String filename, byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment result = GoSlice.allocate(arena);
            GoEnry.GetLanguages(toGoString(arena, filename), toGoByteSlice(arena, content), result);
            return toJavaStringArray(result);
        }
    }

    /**
     * Returns the mime type of the file.
     *
     * @param path     of the file
     * @param language of the file
     * @return mime type
     */
    public static synchronized String getMimeType(String path, String language) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaString(GoEnry.GetMimeType(arena, toGoString(arena, path), toGoString(arena, language)));
        }
    }

    /**
     * Reports whether the given file content is binary or not.
     *
     * @param content of the file
     * @return whether it's binary or not
     */
    public static synchronized boolean isBinary(byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsBinary(toGoByteSlice(arena, content)));
        }
    }

    /**
     * Reports whether the given file or directory is a config file or directory.
     *
     * @param path of the file or directory
     * @return whether it's config or not
     */
    public static synchronized boolean isConfiguration(String path) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsConfiguration(toGoString(arena, path)));
        }
    }

    /**
     * Reports whether the given file or directory it's documentation.
     *
     * @param path of the file or directory. It must not contain its parents and
     *             if it's a directory it must end in a slash e.g. "docs/" or
     *             "foo.json".
     * @return whether it's docs or not
     */
    public static synchronized boolean isDocumentation(String path) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsDocumentation(toGoString(arena, path)));
        }
    }

    /**
     * Reports whether the given file is a dotfile.
     *
     * @param path of the file
     * @return whether it's a dotfile or not
     */
    public static synchronized boolean isDotFile(String path) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsDotFile(toGoString(arena, path)));
        }
    }

    /**
     * Reports whether the given path is an image or not.
     *
     * @param path of the file
     * @return whether it's an image or not
     */
    public static synchronized boolean isImage(String path) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsImage(toGoString(arena, path)));
        }
    }

    /**
     * Reports whether the given path is a vendor path or not.
     *
     * @param path of the file or directory
     * @return whether it's vendor or not
     */
    public static synchronized boolean isVendor(String path) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsVendor(toGoString(arena, path)));
        }
    }

    /**
     * Reports whether the given file is a generated file.
     *
     * @param path of the file
     * @param content of the file
     * @return whether it's autogenerated or not
     */
    public static synchronized boolean isGenerated(String path, byte[] content) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsGenerated(toGoString(arena, path), toGoByteSlice(arena, content)));
        }
    }

    /**
     * Returns a color code for given language.
     *
     * @param language of the file
     * @return color code
     */
    public static synchronized String getColor(String language) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaString(GoEnry.GetColor(arena, toGoString(arena, language)));
        }
    }

    /**
     * Reports whether the given path is a test path or not.
     *
     * @param path of the file or directory
     * @return whether it's test or not
     */
    public static synchronized boolean isTest(String path) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaBool(GoEnry.IsTest(toGoString(arena, path)));
        }
    }

    /**
     * Returns type for given language.
     *
     * @param language of the file
     * @return type (data, programming, markup, prose)
     */
    public static synchronized String getLanguageType(String language) {
        try (Arena arena = Arena.ofConfined()) {
            return toJavaString(GoEnry.GetLanguageType(arena, toGoString(arena, language)));
        }
    }
}
