import itertools
import re
from datetime import date, datetime, time
from typing import Any, Iterable, List, Match, Optional, Tuple, Union, Dict

from notes.structure import Entry, StructuredNotes

# example ### Fri 2021/02/01
DATE_PATTERN = re.compile("### [a-zA-Z]{3} (?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<day>[0-9]{2})")

# example #### title, other words | 12:00/13:00
# example #### title, other words | 12:00
ENTRY_PATTERN = re.compile(
    r"\* (?P<title>.+) \|? (?P<start_hour>[0-9]{2})?:?(?P<start_minutes>[0-9]{2})?/?(?P<end_hour>[0-9]{2})?:?(?P<end_minutes>[0-9]{2})?"
)


class Markdown:
    def parse(self, notes: str) -> StructuredNotes:
        lines = notes.split("\n")
        first_entry_position = self._find_first_entry(lines)
        raw = lines[first_entry_position:]

        parsed = StructuredNotes()
        parsed.top_section = lines[:first_entry_position]

        day: Optional[date] = None
        entry: Optional[Entry] = None
        position = 0

        try:
            while position < len(raw):
                line = raw[position]

                if matched := DATE_PATTERN.match(line):
                    day = self._parse_day(matched)

                    if position + 1 < len(raw) and raw[position + 1] == "":
                        position += 1

                elif (matched := ENTRY_PATTERN.match(line)) and day:
                    title, start, end = self._parse_entry(matched)
                    entry = parsed.add_entry(title=title, timestamp=datetime.combine(day, start))
                    if end:
                        entry.end(datetime.combine(day, end))

                elif entry:
                    entry.append(line)

                else:
                    parsed.legacy.append(line)

                position += 1

        except Exception as e:
            raise ValueError(f"Failed parsing line: {position}, {raw[position]}, {e}")

        return parsed

    def _parse_day(self, matched: Match[Any]) -> date:
        groups = matched.groupdict()
        return datetime(int(groups["year"]), int(groups["month"]), int(groups["day"]))

    def _parse_entry(self, matched: Match[Any]) -> Tuple[str, Optional[time], Optional[time]]:
        groups = matched.groupdict()
        start = self._prase_time(groups, "start")
        end = self._prase_time(groups, "end")

        if not start:
            start = time(0, 0)

        return groups["title"], start, end

    def _prase_time(self, groups: Dict[str, Any], key: str) -> Optional[time]:
        hour = key + "_hour"
        minutes = key + "_minutes"
        if hour in groups and groups[hour]:
            return time(int(groups[hour]), int(groups[minutes]))

    def _find_first_entry(self, lines: List[str]) -> int:
        for i, line in enumerate(lines):
            if DATE_PATTERN.match(line):
                return i

        return len(lines)

    def format(self, notes: StructuredNotes) -> str:
        return "\n".join(
            [
                *notes.top_section,
                self._format_entries(notes.entries),
                *notes.legacy,
            ]
        )

    def _format_entries(self, entries: Iterable[Entry]) -> str:
        sorted_entries = sorted(entries, key=lambda entry: entry.start_time, reverse=True)
        entries_by_day = itertools.groupby(sorted_entries, key=lambda entry: entry.start_time.date())
        return "\n".join(map(lambda day_of_entries: self._format_day_of_entries(*day_of_entries), entries_by_day))

    def _format_day_of_entries(self, day: date, entries: Iterable[Entry]) -> str:
        sorted_entries = sorted(entries, key=lambda entry: entry.start_time)
        formated_entries = map(self._format_entry, sorted_entries)
        return "\n".join(itertools.chain([self._format_day(day)], formated_entries)) + "\n"

    def _format_day(self, timestamp: Union[datetime, date]) -> str:
        return timestamp.strftime("### %a %Y/%m/%d\n")

    def _format_entry(self, entry: Entry) -> str:
        fmt = "%H:%M"
        formatted_entry = f"* {entry.title}"

        if not entry.start_time.time == time(0, 0):
            formatted_entry += f" | {entry.start_time.strftime(fmt)}"

        if entry.end_time:
            formatted_entry += f"/{entry.end_time.strftime(fmt)}"

        return "\n".join([formatted_entry, *entry.content])
