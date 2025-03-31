package events

import (
	"fmt"
	"github.com/ahodieb/papyrus/editor"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func NewEvent(path string) error {
	e := editor.NeovidePrompt
	p := filepath.Join(os.TempDir(), fmt.Sprintf("papyrus-%d.md", time.Now().Unix()))
	if err := e.Open(p, 0); err != nil {
		return err
	}

	//content, err := os.ReadFile(p)
	//if err != nil {
	//	return err
	//}

	//f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	//if err != nil {
	//	return fmt.Errorf("failed to open file %q, %w", path, err)
	//}
	//defer f.Close()
	//
	//event := Event{
	//	Content: string(content),
	//	Time:    time.Now(),
	//}
	//
	//content, err = json.Marshal(event)
	//if err != nil {
	//	return err
	//}
	//
	//if _, err := f.Write(content); err != nil {
	//	return err
	//}

	return nil
}
