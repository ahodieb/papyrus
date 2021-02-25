#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import os
from datetime import datetime
from pathlib import Path
from typing import List, Optional

from notes.manager import NoteManager

_NOTES_FILE_ENV = "JOURNAL_FILE"


def default_command(args: List[str]):
    """
    When papyrus is called with no arguments it runs this command

    This command tries to guess and execute the best action
    1. by default it will create a new entry similar to "papyrus --new"
    2. if the clipboard checking rules match the text in clipboard it adds a new entry with specific structure
    """

    _new_entry(_get_path(), datetime.now(), " ".join(args))


def new_entry_command(args: argparse.Namespace, remaining: List[str]):
    _new_entry(_get_path(args), datetime.now(), " ".join(remaining))


def _new_entry(path: Path, timestamp: datetime, line: str):
    notes = NoteManager(path)
    notes.new_entry(line, timestamp=timestamp)


def _get_path(args: Optional[argparse.Namespace] = None) -> Path:
    path = None

    if args:
        path = args.path
    else:
        path = os.environ[_NOTES_FILE_ENV]

    if not path or Path(path).is_dir():
        raise FileNotFoundError("Path not specified, use --path or set JOURNAL_FILE environment variable")

    return Path(path)
