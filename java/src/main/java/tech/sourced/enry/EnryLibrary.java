package tech.sourced.enry;

import com.sun.jna.Library;
import com.sun.jna.Native;
import com.sun.jna.Pointer;
import com.sun.jna.Platform;
import java.io.File;
import java.io.IOException;
import java.net.URL;
import com.sun.jna.ptr.IntByReference;
import java.util.Collections;

/** Explicit mapping of shared/enry.go. Returned pointers are caller-owned. */
interface EnryLibrary extends Library {
    EnryLibrary INSTANCE = loadBundled();

    static EnryLibrary loadBundled() {
        String resource = "/" + Platform.RESOURCE_PREFIX + "/" + System.mapLibraryName("enry");
        URL location = EnryLibrary.class.getResource(resource);
        if (location == null) throw new UnsatisfiedLinkError("Missing bundled native library: " + resource);
        try {
            File library = Native.extractFromResourcePath(resource, EnryLibrary.class.getClassLoader());
            if (!"file".equals(location.getProtocol())) library.deleteOnExit();
            return Native.load(library.getAbsolutePath(), EnryLibrary.class,
                    Collections.singletonMap(Library.OPTION_STRING_ENCODING, "UTF-8"));
        } catch (IOException error) {
            UnsatisfiedLinkError failure = new UnsatisfiedLinkError("Cannot extract " + resource);
            failure.initCause(error);
            throw failure;
        }
    }

    void FreeCString(Pointer value);
    void FreeStringArray(Pointer value);
    Pointer GetLanguage(String filename, byte[] content, int length);
    Pointer GetLanguageByContentWithSafety(String filename, byte[] content, int length, IntByReference safe);
    Pointer GetLanguageByExtensionWithSafety(String filename, IntByReference safe);
    Pointer GetLanguageByFilenameWithSafety(String filename, IntByReference safe);
    Pointer GetLanguageByEmacsModelineWithSafety(byte[] content, int length, IntByReference safe);
    Pointer GetLanguageByModelineWithSafety(byte[] content, int length, IntByReference safe);
    Pointer GetLanguageByShebangWithSafety(byte[] content, int length, IntByReference safe);
    Pointer GetLanguageByVimModelineWithSafety(byte[] content, int length, IntByReference safe);
    Pointer GetLanguages(String filename, byte[] content, int length);
    Pointer GetLanguageExtensions(String language);
    Pointer GetMimeType(String path, String language);
    Pointer GetColor(String language);
    Pointer GetLanguageType(String language);
    int IsBinary(byte[] content, int length);
    int IsConfiguration(String path);
    int IsDocumentation(String path);
    int IsDotFile(String path);
    int IsImage(String path);
    int IsVendor(String path);
    int IsGenerated(String path, byte[] content, int length);
    int IsTest(String path);
}
