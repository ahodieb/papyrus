from functools import partial

import cli
from cli import commands


def test_parse_default_command():
    args = ["free", "text", "items"]
    command = cli.parse(args)

    assert isinstance(command, partial)
    assert command.func is commands.default_command


def test_parse_empty():
    command = cli.parse([])

    assert isinstance(command, partial)
    assert command.func is commands.default_command


def test_parse_new(tmp_path):
    command = cli.parse(["new", "-p", f"{tmp_path}/notes.txt", "free", "text"])

    assert isinstance(command, partial)
    assert command.func is commands.new_entry_command
