package tech.sourced.enry.internal;

import tech.sourced.enry.NativeLibraryLookup;

import java.io.IOException;
import java.io.InputStream;
import java.lang.foreign.Arena;
import java.lang.foreign.Linker;
import java.lang.foreign.SymbolLookup;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.nio.file.attribute.FileAttribute;
import java.nio.file.attribute.PosixFilePermission;
import java.nio.file.attribute.PosixFilePermissions;
import java.util.Optional;
import java.util.ServiceLoader;
import java.util.Set;

@SuppressWarnings("unused")
final class ChainedLibraryLookup implements NativeLibraryLookup {

    private ChainedLibraryLookup() {}

    static ChainedLibraryLookup INSTANCE = new ChainedLibraryLookup();

    @Override
    public SymbolLookup get(Arena arena) {
        ServiceLoader<NativeLibraryLookup> serviceLoader = ServiceLoader.load(NativeLibraryLookup.class);
        // NOTE: can't use _ because of palantir/palantir-java-format#934
        SymbolLookup lookup = (name) -> Optional.empty();
        for (NativeLibraryLookup libraryLookup : serviceLoader) {
            lookup = lookup.or(libraryLookup.get(arena));
        }

        return lookup
            .or(findResource(arena))
            .or(findLibrary(arena))
            .or(Linker.nativeLinker().defaultLookup());
    }

    private static String libraryPath() {
        String os = System.getProperty("os.name");
        String arch = System.getProperty("os.arch");

        if (os.contains("Linux") && arch.contains("amd64")) {
            return "/linux-x86-64/libenry.so";
        } else if (os.contains("Darwin")) {
            return "/darwin/libenry.dylib";
        }

        return null;
    }

    private static SymbolLookup findResource(Arena arena) {
        String path = libraryPath();
        if (path == null) return SymbolLookup.loaderLookup();

        try (InputStream in = ChainedLibraryLookup.class.getResourceAsStream(path)) {
            if (in == null) return SymbolLookup.loaderLookup();

            Set<PosixFilePermission> perms = PosixFilePermissions.fromString("rw-------");
            FileAttribute<Set<PosixFilePermission>> attribute = PosixFilePermissions.asFileAttribute(perms);

            Path temp = Files.createTempFile("enry-",".tmp", attribute);
            Files.copy(in, temp, StandardCopyOption.REPLACE_EXISTING);
            temp.toFile().deleteOnExit();

            return SymbolLookup.libraryLookup(temp, arena);
        } catch (IOException e) {
            return SymbolLookup.loaderLookup();
        }
    }

    private static SymbolLookup findLibrary(Arena arena) {
        try {
            String library = System.mapLibraryName("enry");
            return SymbolLookup.libraryLookup(library, arena);
        } catch (IllegalArgumentException e) {
            return SymbolLookup.loaderLookup();
        }
    }
}
