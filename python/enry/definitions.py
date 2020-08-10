"""
Python library calling enry Go implementation trough cFFI (API, out-of-line) and Cgo.
"""
from typing import List

from _c_enry import lib
from enry.types import Guess
from enry.utils import transform_types, transform_types_ret_str_slice

GetLanguage = transform_types([str, bytes], str)(lib.GetLanguage)
GetLanguageByContent = transform_types([str, bytes], Guess)(lib.GetLanguageByContent)
GetLanguageByExtension = transform_types([str], Guess)(lib.GetLanguageByExtension)
GetLanguageByFilename = transform_types([str], Guess)(lib.GetLanguageByFilename)
GetLanguageByModeline = transform_types([bytes], Guess)(lib.GetLanguageByModeline)
GetLanguageByShebang = transform_types([bytes], Guess)(lib.GetLanguageByShebang)
GetLanguageByEmacsModeline = transform_types([bytes], Guess)(lib.GetLanguageByEmacsModeline)
GetLanguageByVimModeline = transform_types([bytes], Guess)(lib.GetLanguageByVimModeline)

GetLanguages = transform_types_ret_str_slice([str, bytes])(lib.GetLanguages)
GetLanguageExtensions = transform_types_ret_str_slice([str])(lib.GetLanguageExtensions)

GetMimeType = transform_types([str, str], str)(lib.GetMimeType)
GetColor = transform_types([str], str)(lib.GetColor)

IsVendor = transform_types([str], bool)(lib.IsVendor)
IsGenerated = transform_types([str, bytes], bool)(lib.IsGenerated)
IsBinary = transform_types([bytes], bool)(lib.IsBinary)
IsConfiguration = transform_types([str], bool)(lib.IsConfiguration)
IsDocumentation = transform_types([str], bool)(lib.IsDocumentation)
IsDotFile = transform_types([str], bool)(lib.IsDotFile)
IsImage = transform_types([str], bool)(lib.IsImage)


def get_language(filename: str, content: bytes) -> str:
    """
    Return the language of the given file based on the filename and its contents.

    :param filename: name of the file with the extension
    :param content: array of bytes with the contents of the file (the code)
    :return: the guessed language
    """
    return GetLanguage(filename, content)


def get_language_by_content(filename: str, content: bytes) -> Guess:
    """
    Return detected language by its content.
    If there are more than one possible language, return the first language
    in alphabetical order and safe = False.

    :param filename: path of the file
    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    return GetLanguageByContent(filename, content)


def get_language_by_extension(filename: str) -> Guess:
    """
    Return detected language by the extension of the filename.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param filename: path of the file
    :return: guessed result
    """
    return GetLanguageByExtension(filename)


def get_language_by_filename(filename: str) -> Guess:
    """
    Return detected language by its filename.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param filename: path of the file
    :return: guessed result
    """
    return GetLanguageByFilename(filename)


def get_language_by_modeline(content: bytes) -> Guess:
    """
    Return detected language by its modeline.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    return GetLanguageByModeline(content)


def get_language_by_vim_modeline(content: bytes) -> Guess:
    """
    Return detected language by its vim modeline.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    return GetLanguageByVimModeline(content)


def get_language_by_emacs_modeline(content: bytes) -> Guess:
    """
    Return detected langauge by its emacs modeline.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    return GetLanguageByEmacsModeline(content)


def get_language_by_shebang(content: bytes) -> Guess:
    """
    Return detected langauge by its shebang.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    return GetLanguageByShebang(content)


def get_languages(filename: str, content: bytes) -> List[str]:
    """
    Return all possible languages for the given file.

    :param filename:
    :param content: array of bytes with the contents of the file (the code)
    :return: all possible languages
    """
    return GetLanguages(filename, content)


def get_language_extensions(language: str) -> List[str]:
    """
    Return all the possible extensions for the given language.

    :param language: language to get extensions from
    :return: extensions for given language
    """
    return GetLanguageExtensions(language)


def get_mime_type(path: str, language: str) -> str:
    """
    Return mime type of the file.

    :param path: path of the file
    :param language: language to get mime type from
    :return: mime type
    """
    return GetMimeType(path, language)


def get_color(language: str) -> str:
    """
    Return color code for given language

    :param language:
    :return: color in hex format
    """
    return GetColor(language)


def is_vendor(filename: str) -> bool:
    """
    Return True if given file is a vendor file.

    :param filename: path of the file
    :return: whether it's vendor or not
    """
    return IsVendor(filename)


def is_generated(filename: str, content: bytes) -> bool:
    """
    Return True if given file is a generated file.

    :param filename: path of the file
    :param content: array of bytes with the contents of the file (the code)
    :return: whether it's generated or not
    """
    return IsGenerated(filename, content)


def is_binary(content: bytes) -> bool:
    """
    Return True if given file is a binary file.

    :param content: array of bytes with the contents of the file (the code)
    :return: whether it's binary or not
    """
    return IsBinary(content)


def is_configuration(path: str) -> bool:
    """
    Return True if given file is a configuration file.

    :param path: path of the file
    :return: whether it's a configuration file or not
    """
    return IsConfiguration(path)


def is_documentation(path: str) -> bool:
    """
    Return True if given file is a documentation file.

    :param path: path of the file
    :return: whether it's documentation or not
    """
    return IsDocumentation(path)


def is_dot_file(path: str) -> bool:
    """
    Return True if given file is a dot file.

    :param path: path of the file
    :return: whether it's a dot file or not
    """
    return IsDotFile(path)


def is_image(path: str) -> bool:
    """
    Return True if given file is an image file.

    :param path: path of the file
    :return: whether it's an image or not
    """
    return IsImage(path)
