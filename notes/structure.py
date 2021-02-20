from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


class StructuredNotes:
    def __init__(self):
        self.entries: List[Entry] = []
        self.top_section: List[str] = []
        # needed to handle entries with older formats, this will be removed after implementing legacy format parsers
        self.legacy: List[str] = []

    def add_entry(self, title: str, timestamp: datetime, content: Optional[List[str]] = None) -> "Entry":
        entry = Entry(title=title, start_time=timestamp, content=content if content else [])
        self.entries.append(entry)
        return entry

    def __len__(self) -> int:
        return len(self.top_section) + sum(map(len, self.entries)) + len(self.legacy)


@dataclass
class Entry:
    title: str
    start_time: datetime
    end_time: Optional[datetime] = None
    content: List[str] = field(default_factory=lambda: [])

    def append(self, line: str) -> "Entry":
        self.content.append(line)
        return self

    def end(self, timestamp: datetime):
        self.end_time = timestamp

    def __len__(self) -> int:
        return 1 + len(self.content)