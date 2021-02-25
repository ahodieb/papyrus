import filecmp
from pathlib import Path
from typing import Iterable

from notes.manager import NoteManager


def create_test_file(path: Path, content: Iterable[str]) -> Path:
    file = path / "test-notefile.txt"
    file.write_text("\n".join(content), "utf-8")
    return file


def test_backup(tmp_path: Path):
    expected = create_test_file(tmp_path, ["testing-backup"])
    backup = NoteManager(path=expected).backup()

    assert backup.read_text()
    assert filecmp.cmp(backup, expected)
