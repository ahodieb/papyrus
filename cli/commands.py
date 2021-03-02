#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import os
from datetime import datetime
from pathlib import Path
from typing import List, Optional

from notes.manager import NoteManager

ENV_NOTES_FILE = "JOURNAL_FILE"
ENV_EDITOR = "NOTES_EDITOR"


def default_command(args: List[str]):
    """
    When papyrus is called with no arguments it runs this command

    This command tries to guess and execute the best action
    1. by default it will create a new entry similar to "papyrus --new"
    2. if the clipboard checking rules match the text in clipboard it adds a new entry with specific structure
    """

    notes = NoteManager(_get_path())
    position = notes.new_entry(" ".join(args), timestamp=datetime.now())
    notes.open(editor(), position)


def open_to_latest_entry_command(args: argparse.Namespace, remaining: List[str]):
    notes = NoteManager(_get_path(args))
    notes.open_latest_entry(editor(), timestamp=datetime.now())


def new_entry_command(args: argparse.Namespace, remaining: List[str]):
    notes = NoteManager(_get_path(args))
    notes.new_entry(" ".join(remaining), timestamp=datetime.now())


def editor() -> str:
    return os.environ.get(ENV_EDITOR, os.environ.get("EDITOR", "vim"))


def _get_path(args: Optional[argparse.Namespace] = None) -> Path:
    path = None

    if args:
        path = args.path
    else:
        path = os.environ[ENV_NOTES_FILE]

    if not path or Path(path).is_dir():
        raise FileNotFoundError("Path not specified, use --path or set JOURNAL_FILE environment variable")

    return Path(path)
