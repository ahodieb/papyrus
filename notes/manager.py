import shutil
from datetime import datetime, timedelta
from pathlib import Path
from typing import Tuple

from notes import editors, markdown
from notes.file import NoteFile


class NoteManager:
    def __init__(self, path: Path):
        self.path = path
        self.notes = NoteFile(path)

    def new_entry(self, title: str, timestamp: datetime) -> int:
        """
        add a new entry to notes, and return the line position of the new entry

        :param title str: title for the new entry
        :param timestamp datetime: timestamp for the new entry
        :return int: the line position for the new entry in the notes file
        """

        lines = [markdown.format_entry(title, timestamp)]

        # check if todays date is already in the file
        if not self._find_entry_for_timestamp(timestamp)[1]:
            lines.insert(0, markdown.format_date(timestamp) + "\n")

        position = self._find_latest_position(timestamp)
        self.backup()
        return self.notes.write("\n".join(lines) + "\n", position)

    def open(self, editor: str, position: int):
        """
        open notes in an editor at the specified position

        :param editor str: the editor to use for opening the notes
        :param position int: the position to open the editor at
        """
        editors.open(editor, str(self.path), position)

    def backup(self) -> Path:
        """backup the notes file"""
        backup_dir = self.path.parent / ".backups"
        backup_dir.mkdir(exist_ok=True)
        backup_path = backup_dir / datetime.utcnow().strftime("%Y%m%dT%H%M%S.txt")

        if self.path.exists():
            shutil.copy(self.path, backup_path)

        return backup_path

    def _find_latest_position(self, timestamp: datetime) -> int:
        """
        try to look for the entry from the day before, and append the new entry two lines above it
        if no entry was found (e.g. after the weekend or a vacation) find the latest entry instead

        There is an edge case when there is only one entry in the file for today (e.g. first day of the month)
        so there is no entry for yesterday and an entry from today, then the position should be at the end of the file
        i don't expect to see this case anytime, if it happens i might add some work around
        """

        position, found = self._find_entry_for_timestamp(timestamp - timedelta(days=1))
        if not found:
            position, found = self.notes.find(lambda line: markdown.DATE_PATTERN.match(line))

        return position - 1 if position >= 1 else 0

    def _find_entry_for_timestamp(self, timestamp: datetime) -> Tuple[int, bool]:
        """find the position of a specific timestamp"""
        formatted = markdown.format_date(timestamp)
        return self.notes.find(lambda line: formatted in line)
