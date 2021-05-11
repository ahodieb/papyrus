package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var defaultEditor = Vim // TODO configuration

type EditorOpener interface {
	Open(path string, position int) error
}

type EditorOpenerFunc func(path string, position int) error

func (e EditorOpenerFunc) Open(path string, position int) error {
	return e(path, position)
}

var VSCode EditorOpenerFunc = func(path string, position int) error {
	return open("code", "--g", fmt.Sprintf("%s:%d", path, position))
}

var Vim EditorOpenerFunc = func(path string, position int) error {
	return open("vim", fmt.Sprintf("+%d", position), path)
}

func ByName(name string) EditorOpener {

	n := strings.ToLower(name)

	if n == "vim" {
		return Vim
	}

	if n == "vscode" || n == "code" {
		return VSCode
	}

	return defaultEditor
}

func open(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
