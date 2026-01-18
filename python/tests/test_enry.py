from enry import *

import pytest


@pytest.mark.parametrize("filename,content,language", [
    ("test.py", "import os", "Python"),
    ("", "#!/usr/bin/bash", "Shell"),
    ("test.hs", "", "Haskell"),
])
def test_get_language(filename: str, content: str, language: str):
    # Store encoded content in a variable to prevent early garbage collection
    b_content = content.encode()
    assert get_language(filename, b_content) == language

def test_get_language_by_filename():
    assert get_language_by_filename("pom.xml").language == "Maven POM"


def test_get_language_by_content():
    b_content = "<?php $foo = bar();".encode()
    assert get_language_by_content("test.php", b_content).language == "PHP"

def test_get_language_by_emacs_modeline():
    modeline = "// -*- font:bar;mode:c++ -*-\ntemplate <typename X> class { X i; };"
    b_modeline = modeline.encode()
    assert get_language_by_emacs_modeline(b_modeline).language == "C++"

def test_get_language_by_vim_modeline():
    modeline = "# vim: noexpandtab: ft=javascript"
    b_modeline = modeline.encode()
    assert get_language_by_vim_modeline(b_modeline).language == "JavaScript"

@pytest.mark.parametrize("modeline,language", [
    ("// -*- font:bar;mode:c++ -*-\ntemplate <typename X> class { X i; };", "C++"),
    ("# vim: noexpandtab: ft=javascript", "JavaScript")
])
def test_get_language_by_modeline(modeline: str, language: str):
    b_modeline = modeline.encode()
    assert get_language_by_modeline(b_modeline).language == language

def test_get_language_by_extension():
    assert get_language_by_extension("test.lisp").language == "Common Lisp"


def test_get_language_by_shebang():
    b_shebang = "#!/usr/bin/python3".encode()
    assert get_language_by_shebang(b_shebang).language == "Python"

def test_get_mime_type():
    assert get_mime_type("test.rb", "Ruby") == "text/x-ruby"


def test_is_binary():
    b_content = "println!('Hello world!\n');".encode()
    assert is_binary(b_content) == False

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
    b_content = "import os".encode()
    assert get_languages("test.py", b_content)

def test_get_language_extensions():
    assert get_language_extensions("Python") == [
        ".py", ".cgi", ".fcgi", ".gyp", ".gypi", ".lmi", ".py3", ".pyde",
        ".pyi", ".pyp", ".pyt", ".pyw", ".rpy", ".spec", ".tac",
        ".wsgi", ".xpy"
    ]