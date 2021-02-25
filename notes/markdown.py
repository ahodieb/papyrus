import re
from datetime import date, datetime
from typing import Optional

# example ### Fri 2021/02/01
DATE_PATTERN = re.compile("### [a-zA-Z]{3} (?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<day>[0-9]{2})")

# example #### title, other words | 12:00/13:00
# example #### title, other words | 12:00
ENTRY_PATTERN = re.compile(
    r"\* (?P<title>.+) \|? (?P<start_hour>[0-9]{2})?:?(?P<start_minutes>[0-9]{2})?/?(?P<end_hour>[0-9]{2})?:?(?P<end_minutes>[0-9]{2})?"
)


def format_date(timestamp: datetime) -> str:
    return timestamp.strftime("### %a %Y/%m/%d")


def format_entry(title: str, start_time: datetime, end_time: Optional[datetime] = None) -> str:
    time_format = "%H:%M"
    formatted_entry = f"* {title} | {start_time.strftime(time_format)}"
    if end_time:
        formatted_entry += f"/{end_time.strftime(time_format)}"

    return formatted_entry
