package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const noteExt = "md"

func main() {
	app := cli.NewApp()
	app.Name = "tn"
	app.Usage = "take note"

	app.Commands = []cli.Command{
		{
			Name:    "edit",
			Usage:   "edit a note",
			Aliases: []string{"ed"},
			Action:  editNote,
		},
		{
			Name:    "remove",
			Usage:   "remove a note",
			Aliases: []string{"rm"},
			Action:  removeNote,
		},
		{
			Name:    "list",
			Usage:   "list existing notes",
			Aliases: []string{"ls"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "reverse,r",
				},
				cli.BoolFlag{
					Name: "time,t",
				},
			},
			Action:                 listNotes,
			UseShortOptionHandling: true,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func editNote(ctx *cli.Context) error {
	editor := getEditor()
	note := ctx.Args().First()
	location := noteLocation(note)

	cmd := exec.Command(editor, location)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func removeNote(ctx *cli.Context) error {
	note := ctx.Args().First()
	location := noteLocation(note)

	return os.Remove(location)
}

func listNotes(ctx *cli.Context) error {
	noteDir := getNoteDir()
	noteFiles, err := ioutil.ReadDir(noteDir)
	if err != nil {
		return err
	}

	timeSort := ctx.Bool("time")
	reverse := ctx.Bool("reverse")

	var toSort sort.Interface
	if timeSort {
		toSort = SortByLastModified(noteFiles)
	} else {
		toSort = SortByFileName(noteFiles)
	}

	if reverse {
		toSort = sort.Reverse(toSort)
	}

	sort.Sort(toSort)

	for _, fileInfo := range noteFiles {
		fmt.Println(noteName(fileInfo.Name()))
	}

	return nil
}

type SortByFileName []os.FileInfo

func (f SortByFileName) Len() int           { return len(f) }
func (f SortByFileName) Less(i, j int) bool { return strings.Compare(f[i].Name(), f[j].Name()) < 0 }
func (f SortByFileName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type SortByLastModified []os.FileInfo

func (f SortByLastModified) Len() int           { return len(f) }
func (f SortByLastModified) Less(i, j int) bool { return f[i].ModTime().After(f[j].ModTime()) } // this is backward to make latest first in lists
func (f SortByLastModified) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func removeFileExt(filename string, ext string) string {
	if strings.HasSuffix(filename, ext) {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}

func getAppDir() string {
	dir := os.Getenv("TN_DIR")

	if dir != "" {
		return dir
	}

	dir = os.Getenv("HOME") + "/.tn"

	return dir
}

func getNoteDir() string {
	return getAppDir() + "/notes"
}

func getEditor() string {
	editor := os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}

	editor = "vim"

	return editor
}

func noteLocation(name string) string {
	return getNoteDir() + "/" + name + "." + noteExt
}

func noteName(filename string) string {
	return strings.TrimSuffix(filename, "."+noteExt)
}
