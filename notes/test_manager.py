import filecmp
from datetime import datetime
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


def test_new_entry(tmp_path: Path):
    file = tmp_path / "notes.txt"
    notes = NoteManager(file)

    timestamp = datetime(2021, 2, 1, 10, 35)
    notes.new_entry("Title", timestamp)

    actual = file.read_text("utf-8")

    assert actual == "\n### Mon 2021/02/01\n\n* Title | 10:35\n"
