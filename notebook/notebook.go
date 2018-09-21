package notebook

import (
	"fmt"
	"github.com/brandon-atkinson/take-note/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// note book is a collection of notes
// we don't need to know how they are stored, but we can create, read, update, append, and delete them, referencing them by name

type NoteInfo struct {
	Name         string
	LastModified time.Time
}

func GetAllNoteInfo() ([]NoteInfo, error) {
	fsInfo, err := ioutil.ReadDir(config.NoteDir())
	if err != nil {
		return nil, err
	}

	notesInfo := make([]NoteInfo, len(fsInfo))
	for i, f := range fsInfo {
		notesInfo[i] = NoteInfo{noteName(f.Name()), f.ModTime()}
	}

	return notesInfo, nil
}

func GetNoteInfo(name string) (NoteInfo, error) {
	fsInfo, err := os.Stat(noteLocation(name))
	if err != nil {
		return NoteInfo{}, err
	}

	if fsInfo != nil {
		return NoteInfo{noteName(fsInfo.Name()), fsInfo.ModTime()}, nil
	}

	return NoteInfo{}, nil
}

// create
func StoreNote(name string, content string) error {
	file, err := os.Create(noteLocation(name))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

func NoteExists(name string) (bool, error) {
	info, err := GetNoteInfo(name)
	fmt.Printf("noteinfo for name: %v -> %v", name, info)
	if err != nil {
		return false, err
	}
	return (info != NoteInfo{}), nil
}

// read
func FetchNote(name string) (string, error) {
	bytes, err := ioutil.ReadFile(noteLocation(name))
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// delete
func DeleteNote(name string) error {
	if err := os.Remove(noteLocation(name)); err != nil {
		return err
	}
	return nil
}

func removeFileExt(filename string, ext string) string {
	if strings.HasSuffix(filename, ext) {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}

func noteLocation(name string) string {
	return config.NoteDir() + "/" + name + config.NoteExt()
}

func noteName(filename string) string {
	return strings.TrimSuffix(filename, config.NoteExt())
}
