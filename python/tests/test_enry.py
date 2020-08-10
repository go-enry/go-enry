from enry import *

import pytest


@pytest.mark.parametrize("filename,content,language", [
    ("test.py", "import os", "Python"),
    ("", "#!/usr/bin/bash", "Shell"),
    ("test.hs", "", "Haskell"),
])
def test_get_language(filename: str, content: str, language: str):
    assert get_language(filename, content.encode()) == language


def test_get_language_by_filename():
    assert get_language_by_filename("pom.xml").language == "Maven POM"


def test_get_language_by_content():
    assert get_language_by_content("test.php", "<?php $foo = bar();".encode()).language == "PHP"


def test_get_language_by_emacs_modeline():
    modeline = "// -*- font:bar;mode:c++ -*-\ntemplate <typename X> class { X i; };"
    assert get_language_by_emacs_modeline(modeline.encode()).language == "C++"


def test_get_language_by_vim_modeline():
    modeline = "# vim: noexpandtab: ft=javascript"
    assert get_language_by_vim_modeline(modeline.encode()).language == "JavaScript"


@pytest.mark.parametrize("modeline,language", [
    ("// -*- font:bar;mode:c++ -*-\ntemplate <typename X> class { X i; };", "C++"),
    ("# vim: noexpandtab: ft=javascript", "JavaScript")
])
def test_get_language_by_modeline(modeline: str, language: str):
    assert get_language_by_modeline(modeline.encode()).language == language


def test_get_language_by_extension():
    assert get_language_by_extension("test.lisp").language == "Common Lisp"


def test_get_language_by_shebang():
    assert get_language_by_shebang("#!/usr/bin/python3".encode()).language == "Python"


def test_get_mime_type():
    assert get_mime_type("test.rb", "Ruby") == "text/x-ruby"


def test_is_binary():
    assert is_binary("println!('Hello world!\n');".encode()) == False


@pytest.mark.parametrize("path,is_documentation_actual", [
    ("sss/documentation/", True),
    ("docs/", True),
    ("test/", False),
])
def test_is_documentation(path: str, is_documentation_actual: bool):
    assert is_documentation(path) == is_documentation_actual


@pytest.mark.parametrize("path,is_dot_actual", [
    (".env", True),
    ("something.py", False),
])
def test_is_dot(path: str, is_dot_actual: bool):
    assert is_dot_file(path) == is_dot_actual


@pytest.mark.parametrize("path,is_config_actual", [
    ("configuration.yml", True),
    ("some_code.py", False),
])
def test_is_configuration(path: str, is_config_actual: bool):
    assert is_configuration(path) == is_config_actual


@pytest.mark.parametrize("path,is_image_actual", [
    ("nsfw.jpg", True),
    ("shrek-picture.png", True),
    ("openjdk-1000.parquet", False),
])
def test_is_image(path: str, is_image_actual: bool):
    assert is_image(path) == is_image_actual


def test_get_color():
    assert get_color("Go") == "#00ADD8"


def test_get_languages():
    assert get_languages("test.py", "import os".encode())


def test_get_language_extensions():
    assert get_language_extensions("Python") == [".py", ".cgi", ".fcgi", ".gyp", ".gypi", ".lmi", ".py3", ".pyde",
                                                 ".pyi", ".pyp", ".pyt", ".pyw", ".rpy", ".smk", ".spec", ".tac",
                                                 ".wsgi", ".xpy"]
