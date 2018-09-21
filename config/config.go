package config

import (
	"os"
	"path/filepath"
)

const TN_DIR_ENVVAR = "TN_DIR"

func AppDir() string {
	if dir := os.Getenv(TN_DIR_ENVVAR); dir != "" {
		return dir
	}

	return filepath.Join(os.Getenv("HOME"), ".tn")
}

func NoteDir() string {
	return filepath.Join(AppDir(), "/notes")
}

func Editor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "vi"
}

func NoteExt() string {
	return ".md"
}
