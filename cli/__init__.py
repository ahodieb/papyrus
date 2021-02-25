import argparse
import os
import sys
from functools import partial
from pathlib import Path
from typing import Any, Callable, List, NoReturn, Optional, Text

from cli import commands


class ArgumentParser(argparse.ArgumentParser):
    def error(self, message: Text) -> NoReturn:
        raise argparse.ArgumentError(None, message)


def parser() -> ArgumentParser:
    parser = ArgumentParser()
    parser.add_argument("-p", "--path", help="path to notes file", default=os.environ["JOURNAL_FILE"], type=Path)
    parser.set_defaults(func=commands.default_command)
    sub = parser.add_subparsers()

    new_entry = sub.add_parser("new", help="add new entry to note")
    new_entry.set_defaults(func=commands.new_entry_command)
    return parser


def parse(args: Optional[List[str]] = sys.argv[1:]) -> Callable[[], Any]:
    """
    try to parse args (default sys.argv) if parsing fails, it returns commands.default_command
    """
    try:
        parsed, remaining = parser().parse_known_args(args=args)
        if parsed.func is commands.default_command:
            return partial(commands.default_command, remaining)

        return partial(parsed.func, parsed, remaining)
    except argparse.ArgumentError:
        return partial(commands.default_command, args)
