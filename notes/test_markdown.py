import os
from pathlib import Path
import pytest
from datetime import datetime, timedelta

from notes.formats import Markdown, ENTRY_PATTERN
from notes.structure import StructuredNotes


def test_format_one_entry():
    notes = StructuredNotes()
    timestamp = datetime(2021, 2, 1, 12, 0)

    notes.add_entry(title="title1", timestamp=timestamp).append("line1").append("line2")
    markdown = Markdown()

    actual = markdown.format(notes)
    print(actual)

    assert actual == "### Mon 2021/02/01\n\n* title1 | 12:00\nline1\nline2\n"


def test_format():
    notes = StructuredNotes()
    timestamp = datetime(2021, 2, 1, 12, 0)

    notes.add_entry(title="title1", timestamp=timestamp).append("line1").append("line2")
    notes.add_entry(title="title2", timestamp=timestamp + timedelta(hours=1)).append("line1").append("line2")
    notes.add_entry(title="title3", timestamp=timestamp + timedelta(hours=-1)).append("line1").append("line2")
    notes.add_entry(title="title4", timestamp=timestamp + timedelta(days=1)).append("line1").append("line2")

    expected = "\n".join(
        [
            "### Tue 2021/02/02",
            "",
            "* title4 | 12:00",
            "line1",
            "line2",
            "",
            "### Mon 2021/02/01",
            "",
            "* title3 | 11:00",
            "line1",
            "line2",
            "* title1 | 12:00",
            "line1",
            "line2",
            "* title2 | 13:00",
            "line1",
            "line2",
            "",
        ]
    )

    actual = Markdown().format(notes)

    assert actual == expected


def test_parse_day_only():
    parsed = Markdown().parse("### Mon 2021/02/0")
    assert not parsed.entries


def test_parse_top_section_only():
    parsed = Markdown().parse("Lines\nWithout\nDay\nEntries")
    assert not parsed.entries
    assert parsed.top_section == ["Lines", "Without", "Day", "Entries"]


def test_parse_legacy_entries():
    parsed = Markdown().parse("### Mon 2021/02/01\n* entry 1\n* entry 2\n")
    assert parsed.legacy


def test_round_trip():
    expected = "\n".join(
        [
            "TODO:",
            "* [ ] Top section added",
            "",
            "---",
            "something else",
            "",
            "",
            "### Tue 2021/02/02",
            "",
            "* title4 | 12:00",
            "line1",
            "line2",
            "",
            "### Mon 2021/02/01",
            "",
            "* title3 | 11:00",
            "line1",
            "line2",
            "* title1 | 12:00",
            "line1",
            "line2",
            "* title2 | 13:00",
            "line1",
            "line2",
            "",
        ]
    )

    markdown = Markdown()
    parsed = markdown.parse(expected)
    formated = markdown.format(parsed)

    print(formated)

    assert formated == expected


@pytest.mark.skipif("JOURNAL_FILE" not in os.environ, reason="JOURNAL_FILE is not set")
def test_actual_notes():
    original = Path(os.environ["JOURNAL_FILE"]).read_text("utf-8")

    markdown = Markdown()
    parsed = markdown.parse(original)
    formated = markdown.format(parsed)

    temp_path = "/tmp/test-output"
    Path(temp_path).write_text(formated, "utf-8")
    print(f"code --diff {os.environ['JOURNAL_FILE']} {temp_path}")
    assert formated == original


def test_parse_empty():
    Markdown().parse("")
