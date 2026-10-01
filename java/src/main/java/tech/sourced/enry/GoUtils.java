package tech.sourced.enry;

import com.sun.jna.Pointer;

/** Conversion and ownership rules for the shared library's C ABI. */
final class GoUtils {
    private GoUtils() {}

    static String cString(String value) {
        if (value == null) return "";
        if (value.indexOf('\0') >= 0) {
            throw new IllegalArgumentException("C string arguments cannot contain NUL");
        }
        return value;
    }

    static int length(byte[] value) {
        return value == null ? 0 : value.length;
    }

    static String toJavaString(Pointer value) {
        if (value == null) return "";
        try {
            return value.getString(0, "UTF-8");
        } finally {
            EnryLibrary.INSTANCE.FreeCString(value);
        }
    }

    static String[] toJavaStringArray(Pointer value) {
        if (value == null) return new String[0];
        try {
            return value.getStringArray(0, "UTF-8");
        } finally {
            EnryLibrary.INSTANCE.FreeStringArray(value);
        }
    }

    static boolean toJavaBool(int value) {
        return value != 0;
    }
}
