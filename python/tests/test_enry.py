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


@pytest.mark.parametrize("filename,content,candidates", [
    ('test.py', b'print("Hello World")', None),
    ('test.py', b'print("Hello World")', []),
])
def test_get_languages_by_filename_with_empty_or_none_candidates(filename: str, content: bytes, candidates):
    result = get_languages_by_filename(filename, content, candidates)
    assert isinstance(result, list)
    assert result == []


def test_get_languages_by_filename_with_valid_candidates():
    # Use a known filename that GetLanguagesByFilename recognizes
    result = get_languages_by_filename('pom.xml', b'<xml></xml>', ['Maven POM', 'Ruby'])
    assert isinstance(result, list)
    # Should return Maven POM since 'pom.xml' is a known filename
    assert 'Maven POM' in result


@pytest.mark.parametrize("filename,content,candidates,expected", [
    ('test.py', b'#!/usr/bin/env python', None, ['Python']),
    ('test.py', b'#!/usr/bin/env python', [], []),
])
def test_get_languages_by_shebang_with_empty_or_none_candidates(filename: str, content: bytes, candidates, expected):
    result = get_languages_by_shebang(filename, content, candidates)
    assert isinstance(result, list)
    assert result == expected


def test_get_languages_by_shebang_with_valid_candidates():
    result = get_languages_by_shebang('test.sh', b'#!/usr/bin/env python', ['Python', 'Shell', 'Bash'])
    assert isinstance(result, list)
    # Should return Python since the shebang is python
    assert 'Python' in result

@pytest.mark.parametrize("filename,safe", [("test.py", True), ("test.h", False), ("test.unknownextension", False)])
def test_extension_ambiguity(filename, safe):
    assert get_language_by_extension(filename).safe is safe


def test_emacs_modeline_does_not_accept_vim_modeline():
    content = b"# vim: noexpandtab: ft=javascript"
    assert get_language_by_emacs_modeline(content).language == ""
    assert get_language_by_vim_modeline(content).language == "JavaScript"


def test_embedded_nul_content():
    assert is_binary(b"text\x00more text") is True

@pytest.mark.parametrize('detector,filename,content,language', [
    (get_languages_by_filename, 'pom.xml', b'', 'Maven POM'),
    (get_languages_by_shebang, 'script', b'#!/usr/bin/python', 'Python'),
])
@pytest.mark.parametrize('filter_kind', ['none', 'empty', 'excluded', 'included'])
def test_candidate_filters(detector, filename, content, language, filter_kind):
    candidates = {'none': None, 'empty': [], 'excluded': ['Ruby'],
                  'included': ['Ruby', language]}[filter_kind]
    expected = [language] if filter_kind in ('none', 'included') else []
    assert detector(filename, content, candidates) == expected

@pytest.mark.parametrize('options', ['-u PYTHONHOME', '--unset PYTHONHOME', '-C /tmp'])
def test_shebang_env_option_operands(options):
    content = f'#!/usr/bin/env -S {options} python3'.encode()
    guess = get_language_by_shebang(content)
    assert guess.language == 'Python'
    assert guess.safe is True
    assert get_languages_by_shebang('script', content, ['Python']) == ['Python']

@pytest.mark.parametrize('args', ['-Spython3 -u', '--split-string=python3 -u', '-S "python3" -u'])
def test_shebang_split_string(args):
    content = ('#!/usr/bin/env ' + args).encode()
    assert get_language_by_shebang(content).language == 'Python'
    assert get_language('script', content) == 'Python'

@pytest.mark.parametrize('path,mime', [
    ('photo.jpg', 'image/jpeg'), ('photo.JPG', 'image/jpeg'),
    ('photo.jpeg', 'image/jpeg'), ('photo.JPEG', 'image/jpeg'),
    ('photo.PNG', 'image/png'), ('photo.GIF', 'image/gif'),
])
def test_image_mime_types(path, mime):
    assert is_image(path)
    assert get_mime_type(path, '') == mime

@pytest.mark.parametrize('content', [
    b'// Generated by CoffeeScript 2.7.0\nconsole.log("hello");\n',
    b'(function() {\n  var __bind = function(fn, me) {};\n}).call(this);\n',
])
def test_generated_coffeescript(content):
    assert is_generated('sample.js', content)

@pytest.mark.parametrize('ending', [b'', b'\n'])
def test_generated_vcr_footer(ending):
    assert is_generated('cassette.yml', b'---\nhttp_interactions: []\nrecorded_with: VCR 6.0.0' + ending)


def test_generated_footer_does_not_drop_final_character():
    assert not is_generated('cassette.yml', b'recorded_with: VCR 6.0.0\nx')
    assert not is_generated('code.js', b'//# sourceMappingURL=code.js.map\nx\ny')
