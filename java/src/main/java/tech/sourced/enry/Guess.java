package tech.sourced.enry;

/**
 * Guess denotes a language detection result of which enry can be
 * completely sure or not.
 */
public class Guess {
    /**
     * Result is the resultant language of the detection.
     */
    public String result;

    /**
     * Sure indicates whether the enry was completely sure the language is
     * the correct one or it might not be.
     */
    public boolean sure;

    public Guess(String result, boolean sure) {
        this.result = result;
        this.sure = sure;
    }
}
