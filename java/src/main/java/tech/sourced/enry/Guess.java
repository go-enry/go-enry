package tech.sourced.enry;

/**
 * Guess denotes a language detection result of which enry can be
 * completely sure or not.
 */
public class Guess {
    /**
     * The resultant language of the detection.
     */
    public String language;

    /**
     * Indicates whether the enry was completely sure the language is
     * the correct one or it might not be.
     */
    public boolean safe;

    public Guess(String language, boolean safe) {
        this.language = language;
        this.safe = safe;
    }
}
