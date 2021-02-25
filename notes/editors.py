import subprocess


def vscode(path: str, position: int):
    subprocess.call(["code", "--g", f"{path}:{position}"])


def vim(path: str, position: int):
    subprocess.call(["vim", f"+{position}", path])


EDITORS = {
    "vscode": vscode,
    "code": vscode,
    "vim": vim,
}


def open(editor: str, path: str, position: int):
    if editor in EDITORS:
        EDITORS[editor](path, position)
