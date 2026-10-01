package tech.sourced.enry;

import org.junit.Test;
import java.util.HashSet;
import static org.junit.Assert.*;

public class GuessTest {
    @Test
    public void equalResultsAreInterchangeableInHashCollections() {
        Guess first = new Guess("Python", true);
        Guess second = new Guess("Python", true);
        assertEquals(first, second);
        assertEquals(first.hashCode(), second.hashCode());
        HashSet<Guess> results = new HashSet<>();
        results.add(first);
        assertTrue(results.contains(second));
        assertNotEquals(first, new Guess("Python", false));
        assertNotEquals(first, new Guess("Go", true));
        assertNotEquals(first, null);
        assertNotEquals(first, "Python");
    }

    @Test
    public void nullLanguageHasConsistentValueSemantics() {
        assertEquals(new Guess(null, false), new Guess(null, false));
        assertEquals(new Guess(null, false).hashCode(), new Guess(null, false).hashCode());
        assertNotEquals(new Guess(null, false), new Guess("", false));
    }
}
