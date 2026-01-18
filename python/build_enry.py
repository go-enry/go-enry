import os
from pathlib import Path
from cffi import FFI

ffibuilder = FFI()

# We manually define the "Contract". This is clean, standard C that CFFI 
# can parse reliably on every platform without needing the .h file.
ffibuilder.cdef("""
    // Memory Management
    void FreeCString(char* str);
    void FreeStringArray(char** ptr);

    // Singular API
    char* GetLanguage(char* filename, char* content, int length);
    char* GetLanguageByContent(char* filename, char* content, int length);
    char* GetLanguageByEmacsModeline(char* content, int length);
    char* GetLanguageByExtension(char* filename);
    char* GetLanguageByFilename(char* filename);
    char* GetLanguageByModeline(char* content, int length);
    char* GetLanguageByShebang(char* content, int length);
    char* GetLanguageByVimModeline(char* content, int length);
    char* GetMimeType(char* path, char* language);
    char* GetColor(char* language);
    char* GetLanguageType(char* language);

    // Boolean API (using int for stability across platforms)
    int IsBinary(char* content, int length);
    int IsConfiguration(char* path);
    int IsDocumentation(char* path);
    int IsDotFile(char* path);
    int IsImage(char* path);
    int IsVendor(char* path);
    int IsGenerated(char* path, char* content, int length);
    int IsTest(char* path);

    // Plural API (returning char** arrays)
    char** GetLanguages(char* filename, char* content, int length);
    char** GetLanguageExtensions(char* language);
    char** GetLanguagesByContent(char* filename, char* content, int length, char** candidates);
    char** GetLanguagesByEmacsModeline(char* filename, char* content, int length, char** candidates);
    char** GetLanguagesByExtension(char* filename, char* content, int length, char** candidates);
    char** GetLanguagesByFilename(char* filename, char* content, int length, char** candidates);
    char** GetLanguagesByModeline(char* filename, char* content, int length, char** candidates);
    char** GetLanguagesByShebang(char* filename, char* content, int length, char** candidates);
    char** GetLanguagesByVimModeline(char* filename, char* content, int length, char** candidates);
""")

# This creates a module named _c_enry inside the enry package.
# We use 'None' because we are in ABI mode (dlopen) rather than API mode.
ffibuilder.set_source("enry._c_enry", None)

if __name__ == "__main__":
    # Get the directory where this script lives (the 'python/' folder)
    base_dir = Path(__file__).resolve().parent
    
    # Target the 'enry/' subfolder
    target_dir = base_dir / "enry"
    
    # Ensure the folder exists
    target_dir.mkdir(parents=True, exist_ok=True)
    
    # Compile the _c_enry.py module into the target folder
    # In ABI mode, CFFI ignores 'tmpdir' for the final module placement 
    # and uses the name in set_source relative to the current working directory.
    # To be safe, we change the directory or use absolute paths.
    os.chdir(base_dir) 
    
    ffibuilder.compile(verbose=True)
    print(f"Generated CFFI module in {target_dir}")