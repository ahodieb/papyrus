package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var VSCode OpenerFunc = func(path string, position int) error {
	return open("code", "--g", fmt.Sprintf("%s:%d", path, position))
}

var Vim OpenerFunc = func(path string, position int) error {
	return open("vim", fmt.Sprintf("+%d", position), path)
}

var Neovim OpenerFunc = func(path string, position int) error {
	return open("nvim", fmt.Sprintf("+%d", position), path)
}

var Neovide OpenerFunc = func(path string, position int) error {
	return open("neovide", "--fork", fmt.Sprintf("+%d", position), path)
}

var DefaultEditor = Vim

type Opener interface {
	Open(path string, position int) error
}

type OpenerFunc func(path string, position int) error

func (e OpenerFunc) Open(path string, position int) error {
	return e(path, position)
}

func ByName(name string) Opener {
	switch strings.ToLower(name) {
	case "vim":
		return Vim
	case "nvim":
		return Neovim
	case "vscode":
		return VSCode
	case "neovide":
		return Neovide
	default:
		return DefaultEditor
	}
}

func open(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	//out, err := cmd.CombinedOutput()
	//fmt.Println(string(out))
	return cmd.Run()
}
