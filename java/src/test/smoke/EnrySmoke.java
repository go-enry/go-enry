import tech.sourced.enry.Enry;

class EnrySmoke {
    public static void main(String[] args) {
        if (!"Python".equals(Enry.getLanguage("example.py", null))) {
            throw new AssertionError("Packaged native detector did not load");
        }
        if (Enry.getLanguageByExtension("example.h").safe) {
            throw new AssertionError("Ambiguity information was lost");
        }
        System.out.println("Packaged Java/native binding smoke test passed");
    }
}
