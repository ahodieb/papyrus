from pathlib import Path
from typing import Any, Callable, Iterable, Tuple


class NoteFile:
    """Utility class providing read and write operations to note files"""

    def __init__(self, path: Path):
        """
        :param path Path: path of the file, if it does not exist an new file will be created
        """

        self.path: Path = path
        self.path.touch(exist_ok=True)

    def write(self, content: str, position: int = 0) -> int:
        """
        Add content to journal file
        :param content str: content to add to the file
        :param position int (default 0): optional line number to instert the content at
        :return position after the written content
        """

        if not content.endswith("\n"):
            content += "\n"

        with open(self.path, "r+") as file:
            lines = file.readlines()
            lines.insert(position, content)

            file.seek(0)
            for line in lines:
                file.write(line)
            file.flush()

        return position + len(content.split("\n")) - 1

    def enumerate(self) -> Iterable[Tuple[int, str]]:
        """read the file line by line, and include line number in the tuple"""

        file = open(self.path, "r")
        line_number = 0

        for line in file:
            yield line_number, line.strip()
            line_number += 1

        file.close()

    def read(self) -> Iterable[str]:
        """iterates over all lines in the file"""

        return (line for _, line in self.enumerate())

    def find(self, matcher: Callable[[str], Any]) -> Tuple[int, bool]:
        """
        find the first position that matches the matcher function,
        :param matcher: function that returns a bool given a line to match
               example: find(lambda line: "abc" in line)

        :return Tuple[line_number: int, found: bool]
        """

        filtered = ((line_number, True) for line_number, _ in self.find_all(matcher))
        return next(filtered, (0, False))

    def find_all(self, matcher: Callable[[str], Any]) -> Iterable[Tuple[int, str]]:
        """
        find all lines matching the matcher function
        :param matcher: function that returns a bool given a line to match
               example: find(lambda line: "abc" in line)
        """

        return filter(lambda x: matcher(x[1]), self.enumerate())
