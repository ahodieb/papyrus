from datetime import datetime, timedelta

from notes.structure import Entry, StructuredNotes


def test_add_entry():
    notes = StructuredNotes()
    timestamp = datetime(2021, 2, 1, 12, 0)

    notes.add_entry("first", timestamp=timestamp)
    notes.add_entry("second", timestamp=timestamp + timedelta(hours=1))
    notes.add_entry("third", timestamp=timestamp + timedelta(hours=3))

    assert sorted([entry.title for entry in notes.entries]) == ["first", "second", "third"]


def test_entry():
    timestamp = datetime(2021, 2, 1, 12, 0)
    entry = Entry(title="title", start_time=timestamp)

    entry.append("line1").append("line2").end(timestamp + timedelta(hours=1))

    assert entry.content == ["line1", "line2"]
    assert entry.end_time == timestamp + timedelta(hours=1)