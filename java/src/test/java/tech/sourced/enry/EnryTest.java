package tech.sourced.enry;

import org.junit.Test;

import static org.junit.Assert.*;

public class EnryTest {

    @Test
    public void getLanguage() {
        String code = "<?php $foo = bar();";
        assertEquals("PHP", Enry.getLanguage("foobar.php", code.getBytes()));
    }

    @Test
    public void getLanguageWithNullContent() {
        assertEquals("Python", Enry.getLanguage("foo.py",  null));
    }

    @Test
    public void getLanguageWithEmptyContent() {
        assertEquals("Go", Enry.getLanguage("baz.go",  "".getBytes()));
        assertEquals("Go", Enry.getLanguage("baz.go",  null));
    }

    @Test
    public void getLanguageWithNullFilename() {
        byte[] content = "#!/usr/bin/env python".getBytes();
        assertEquals("Python", Enry.getLanguage(null, content));
    }

    @Test
    public void getLanguageByContent() {
        String code = "<?php $foo = bar();";
        assertGuess(
                "PHP",
                true,
                Enry.getLanguageByContent("foo.php", code.getBytes())
        );
    }

    @Test
    public void getLanguageByFilename() {
        assertGuess(
                "Maven POM",
                true,
                Enry.getLanguageByFilename("pom.xml")
        );
    }

    @Test
    public void getLanguageByEmacsModeline() {
        String code = "// -*- font:bar;mode:c++ -*-\n" +
                "template <typename X> class { X i; };";
        assertGuess(
                "C++",
                true,
                Enry.getLanguageByEmacsModeline(code.getBytes())
        );
    }

    @Test
    public void getLanguageByExtension() {
        assertGuess(
                "Ruby",
                true,
                Enry.getLanguageByExtension("foo.rb")
        );
    }

    @Test
    public void getLanguageByShebang() {
        String code = "#!/usr/bin/env python";
        assertGuess(
                "Python",
                true,
                Enry.getLanguageByShebang(code.getBytes())
        );
    }

    @Test
    public void getLanguageByModeline() {
        String code = "// -*- font:bar;mode:c++ -*-\n" +
                "template <typename X> class { X i; };";
        assertGuess(
                "C++",
                true,
                Enry.getLanguageByModeline(code.getBytes())
        );

        code = "# vim: noexpandtab: ft=javascript";
        assertGuess(
                "JavaScript",
                true,
                Enry.getLanguageByModeline(code.getBytes())
        );
    }

    @Test
    public void getLanguageByVimModeline() {
        String code = "# vim: noexpandtab: ft=javascript";
        assertGuess(
                "JavaScript",
                true,
                Enry.getLanguageByVimModeline(code.getBytes())
        );
    }

    @Test
    public void getLanguageExtensions() {
        String[] exts = Enry.getLanguageExtensions("Go");
        String[] expected = {".go"};
        assertArrayEquals(expected, exts);
    }

    @Test
    public void getLanguages() {
        String code = "#include <stdio.h>" +
                "" +
                "extern int foo(void *bar);";

        String[] result = Enry.getLanguages("foo.h", code.getBytes());
        String[] expected = {"C"};
        assertArrayEquals(expected, result);
    }

    @Test
    public void getMimeType() {
        assertEquals(
                "text/x-ruby",
                Enry.getMimeType("foo.rb", "Ruby")
        );
    }

    @Test
    public void isBinary() {
        assertFalse(Enry.isBinary("hello = 'world'".getBytes()));
    }

    @Test
    public void isConfiguration() {
        assertTrue(Enry.isConfiguration("config.yml"));
        assertFalse(Enry.isConfiguration("FooServiceProviderImplementorFactory.java"));
    }

    @Test
    public void isDocumentation() {
        assertTrue(Enry.isDocumentation("docs/"));
        assertFalse(Enry.isDocumentation("src/"));
    }

    @Test
    public void isDotFile() {
        assertTrue(Enry.isDotFile(".env"));
        assertFalse(Enry.isDotFile("config.json"));
    }

    @Test
    public void isImage() {
        assertTrue(Enry.isImage("yup.jpg"));
        assertFalse(Enry.isImage("nope.go"));
    }

    @Test
    public void getColor() {
        assertEquals(
                "#00ADD8",
                Enry.getColor("Go")
        );
    }

    @Test
    public void isTest() {
        assertTrue(Enry.isTest("test_foo.py"));
        assertTrue(Enry.isTest("test/java/tech/sourced/enry/EnryTest.java"));
        assertFalse(Enry.isTest("foo.py"));
        assertFalse(Enry.isTest("src/java/foo.java"));
    }

    @Test
    public void getLanguageType() {
        assertEquals("data", Enry.getLanguageType("CSV"));
    }

    @Test
    public void binaryContentPreservesEmbeddedNul() {
        assertTrue(Enry.isBinary(new byte[] {'a', 0, 'b'}));
        assertFalse(Enry.isBinary(null));
        assertFalse(Enry.isBinary(new byte[0]));
    }

    @Test
    public void emptyAndUnknownResults() {
        assertGuess("", false, Enry.getLanguageByExtension("unknown.zzzzzz"));
        assertArrayEquals(new String[0], Enry.getLanguageExtensions("not a language"));
        assertEquals("Python", Enry.getLanguage("例.py", new byte[0]));
        assertEquals("", Enry.getLanguage(null, null));
    }

    @Test
    public void ambiguousExtensionPreservesSafety() {
        Guess guess = Enry.getLanguageByExtension("sample.h");
        assertFalse(guess.safe);
        assertFalse(guess.language.isEmpty());
    }

    @Test(expected = IllegalArgumentException.class)
    public void rejectsTruncatedCStringArguments() {
        Enry.getLanguage("sample.py\0.txt", null);
    }

    @Test
    public void generatedAndVendorPredicates() {
        assertTrue(Enry.isVendor("node_modules/example/index.js"));
        assertTrue(Enry.isGenerated("example.go",
                "// Code generated by test. DO NOT EDIT.\npackage test\n".getBytes()));
    }

    @Test
    public void repeatedOwnedResultsRemainUsable() {
        for (int i = 0; i < 10000; i++) {
            assertEquals("Python", Enry.getLanguage("sample.py", null));
            assertArrayEquals(new String[] {".go"}, Enry.getLanguageExtensions("Go"));
            assertGuess("Ruby", true, Enry.getLanguageByExtension("sample.rb"));
        }
    }

    void assertGuess(String language, boolean safe, Guess guess) {
        assertEquals(language, guess.language);
        assertEquals(safe, guess.safe);
    }

}
