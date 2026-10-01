from pathlib import Path
from types import SimpleNamespace

import pytest
import enry.definitions as definitions


def test_missing_library_does_not_load_from_working_directory(tmp_path, monkeypatch):
    package = tmp_path / 'install' / 'enry'
    package.mkdir(parents=True)
    untrusted = tmp_path / 'untrusted'
    untrusted.mkdir()
    for name in ('libenry.dylib', 'libenry.so', 'libenry.dll'):
        (untrusted / name).write_bytes(b'not a library')
    monkeypatch.chdir(untrusted)
    monkeypatch.setattr(definitions, '__file__', str(package / 'definitions.py'))
    attempted = []

    def fake_dlopen(path):
        attempted.append(Path(path))
        raise OSError('missing library')

    monkeypatch.setattr(definitions, 'ffi', SimpleNamespace(dlopen=fake_dlopen))
    with pytest.raises(ImportError):
        definitions._load_library()
    assert attempted
    assert all(path.parent != untrusted for path in attempted)
