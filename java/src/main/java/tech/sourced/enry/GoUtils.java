package tech.sourced.enry;

import tech.sourced.enry.internal.GoEnry;
import tech.sourced.enry.internal.GoSlice;
import tech.sourced.enry.internal.GoString;

import java.lang.foreign.Arena;
import java.lang.foreign.MemorySegment;
import java.lang.foreign.ValueLayout;
import java.nio.charset.StandardCharsets;

class GoUtils {

    static MemorySegment toGoString(Arena arena, String str) {
        if (str == null) str = "";

        byte[] bytes = str.getBytes(StandardCharsets.UTF_8);
        MemorySegment p = arena.allocateFrom(ValueLayout.JAVA_BYTE, bytes);

        MemorySegment goString = GoString.allocate(arena);
        GoString.p(goString, p);
        GoString.n(goString, bytes.length);

        return goString;
    }

    static String toJavaString(MemorySegment str) {
        long n = GoString.n(str);
        if (n == 0) return "";

        MemorySegment p = GoString.p(str);
        byte[] bytes = p.asSlice(0, n).toArray(GoEnry.C_CHAR);

        return new String(bytes, StandardCharsets.UTF_8);
    }

    static String[] toJavaStringArray(MemorySegment slice) {
        long len = GoSlice.len(slice);
        MemorySegment data = GoSlice.data(slice);

        String[] result = new String[(int) len];
        for (int i = 0; i < len; i++) {
            MemorySegment s = data.get(GoEnry.C_POINTER, i * GoEnry.C_POINTER.byteSize());
            result[i] = s.getString(0);
        }
        return result;
    }

    static MemorySegment toGoByteSlice(Arena arena, byte[] bytes) {
        if (bytes == null) bytes = new byte[0];

        MemorySegment data = arena.allocateFrom(ValueLayout.JAVA_BYTE, bytes);

        MemorySegment goSlice = GoSlice.allocate(arena);
        GoSlice.data(goSlice, data);
        GoSlice.len(goSlice, bytes.length);
        GoSlice.cap(goSlice, bytes.length);

        return goSlice;
    }

    static boolean toJavaBool(byte goBool) {
        return goBool == 1;
    }
}
