package tech.sourced.enry;

import java.lang.foreign.Arena;
import java.lang.foreign.SymbolLookup;

/**
 * An interface implemented by clients that wish to customize the {@link SymbolLookup}
 * used for the go-enry native library. Implementations must be registered
 * by listing their fully qualified class name in a resource file named
 * {@code META-INF/services/tech.sourced.enry.NativeLibraryLookup}.
 *
 * @since 2.9.4
 * @see java.util.ServiceLoader
 */
@FunctionalInterface
public interface NativeLibraryLookup {
    /**
     * Get the {@link SymbolLookup} to be used for the go-enry native library.
     *
     * @param arena The arena that will manage the native memory.
     * @since 2.9.4
     */
    SymbolLookup get(Arena arena);
}
