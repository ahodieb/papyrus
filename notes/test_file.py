from pathlib import Path
from typing import Iterable

from notes.file import NoteFile


def create_test_file(path: Path, content: Iterable[str]) -> Path:
    file = path / "test-notefile.txt"
    file.write_text("\n".join(content), "utf-8")
    return file


def test_write(tmp_path: Path):
    content = (
        "line 1",
        "line 2",
    )

    notes = NoteFile(path=create_test_file(tmp_path, content))
    notes.write("line 3", 1)

    assert notes.path.read_text("utf-8") == "line 1\nline 3\nline 2"


def test_write_scenario(tmp_path: Path):
    file = NoteFile(tmp_path / "test-notefile.txt")

    file.write("First entry")
    position = file.write("Second entry")
    position = file.write("1", position=position)
    position = file.write("2", position=position)

    assert list(file.read()) == ["Second entry", "1", "2", "First entry"]


def test_enumerate(tmp_path: Path):
    content = (
        "line 1",
        "line 2",
        "line 3",
    )

    notes = NoteFile(path=create_test_file(tmp_path, content))

    for i, line in notes.enumerate():
        assert line == content[i]


def test_read(tmp_path: Path):
    content = [
        "line 1",
        "line 2",
        "line 3",
    ]

    notes = NoteFile(path=create_test_file(tmp_path, content))
    assert list(notes.read()) == content


def test_find(tmp_path: Path):
    content = (
        "line 1",
        "line 2",
        "line 3",
    )

    notes = NoteFile(path=create_test_file(tmp_path, content))
    assert notes.find(lambda line: line == "line 2") == (1, True)


def test_find_not_found(tmp_path: Path):
    notes = NoteFile(path=create_test_file(tmp_path, []))
    assert notes.find(lambda line: False) == (0, False)


def test_find_all(tmp_path: Path):
    content = ("word", "something else", "words")

    notes = NoteFile(path=create_test_file(tmp_path, content))
    assert list(notes.find_all(lambda line: line.startswith("word"))) == [(0, "word"), (2, "words")]
