package tech.sourced.enry;

import com.sun.jna.Memory;
import com.sun.jna.Pointer;
import tech.sourced.enry.nativelib.GoSlice;
import tech.sourced.enry.nativelib.GoString;

import java.io.UnsupportedEncodingException;

class GoUtils {

    static GoString.ByValue toGoString(String str) {
        byte[] bytes;
        try {
            bytes = str.getBytes("utf-8");
        } catch (UnsupportedEncodingException e) {
            bytes = str.getBytes();
        }

        GoString.ByValue val = new GoString.ByValue();
        val.n = bytes.length;
        Pointer ptr = new Memory(bytes.length);
        ptr.write(0, bytes, 0, bytes.length);
        val.p = ptr;
        return val;
    }

    static String toJavaString(GoString str) {
        if (str.n == 0) {
            return "";
        }

        byte[] bytes = new byte[(int) str.n];
        str.p.read(0, bytes, 0, (int) str.n);
        try {
            return new String(bytes, "utf-8");
        } catch (UnsupportedEncodingException e) {
            throw new RuntimeException("utf-8 encoding is not supported");
        }
    }

    static String[] toJavaStringArray(GoSlice slice) {
        String[] result = new String[(int) slice.len];
        Pointer[] ptrArr = slice.data.getPointerArray(0, (int) slice.len);
        for (int i = 0; i < (int) slice.len; i++) {
            result[i] = ptrArr[i].getString(0);
        }
        return result;
    }

    static GoSlice.ByValue toGoByteSlice(byte[] bytes) {
        return sliceFromPtr(bytes.length, ptrFromBytes(bytes));
    }

    static GoSlice.ByValue sliceFromPtr(int len, Pointer ptr) {
        GoSlice.ByValue val = new GoSlice.ByValue();
        val.cap = len;
        val.len = len;
        val.data = ptr;
        return val;
    }

    static Pointer ptrFromBytes(byte[] bytes) {
        Pointer ptr = new Memory(bytes.length);
        ptr.write(0, bytes, 0, bytes.length);
        return ptr;
    }

    static boolean toJavaBool(byte goBool) {
        return goBool == 1;
    }

}
