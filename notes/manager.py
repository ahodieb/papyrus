import shutil
from datetime import datetime, timedelta
from pathlib import Path
from typing import Tuple

from notes import markdown
from notes.structure import StructuredNotes


class NoteManager:
    def __init__(self, path: Path):
        self.path = path
        self.path.touch(exist_ok=True)
        self._parser = Markdown()
        self._formater = self._parser

    def new_entry(self, title: str, timestamp: datetime) -> int:
        notes = self._parser.parse(self.path.read_text("utf-8"))
        notes.add_entry(title=title, timestamp=timestamp)
        self.write(notes)

    def write(self, notes: StructuredNotes):
        self.backup()
        self.path.write_text(self._formater.format(notes), "utf-8")

    def backup(self) -> Path:
        backup_dir = self.path.parent / ".backups"
        backup_dir.mkdir(exist_ok=True)
        backup_path = backup_dir / datetime.utcnow().strftime("%Y%m%dT%H%M%S.txt")

        if self.path.exists():
            shutil.copy(self.path, backup_path)

        return backup_path

    def _find_latest_position(self, timestamp: datetime) -> int:
        """
        Ported the logic from my original bash scripts

        Try to look for the entry from the day before, and append the new entry two lines above it
        if no entry was found (e.g. after the weekend or a vacation) find the latest entry instead

        There is an edge case when there is only one entry in the file for today (e.g. first day of the month)
        so there is no entry for yesterday and an entry from today, then the position should be at the end of the file
        i don't expect to see this case anytime, if it happens i might add some work around
        """

        position, found = self._find_entry_for_timestamp(timestamp - timedelta(days=1))
        if not found:
            position, found = self.notes.find(lambda line: markdown.DATE_PATTERN.match(line))

        return position - 2 if position >= 2 else 0

    def _find_entry_for_timestamp(self, timestamp: datetime) -> Tuple[int, bool]:
        formatted = markdown.format_date(timestamp)
        return self.notes.find(lambda line: formatted in line)
