package main

import (
	"fmt"
	"github.com/brandon-atkinson/take-note/config"
	"github.com/brandon-atkinson/take-note/notebook"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "tn"
	app.Usage = "take note"
	app.EnableBashCompletion = true
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
			Aliases: []string{"rem", "rm"},
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
	editor := config.Editor()
	note := ctx.Args().First()

	tempFile, err := ioutil.TempFile("", note+"*"+config.NoteExt())
	if err != nil {
		return err
	}
	tempFileName := tempFile.Name()

	noteExists, err := notebook.NoteExists(note)
	fmt.Printf("noteExists? note %v %v\n", note, noteExists)
	if err != nil {
		return err
	}

	if noteExists {
		content, err := notebook.FetchNote(note)
		fmt.Printf("content? %v\n", content)
		if err != nil {
			return err
		}

		_, err = tempFile.WriteString(content)
		if err != nil {
			return err
		}
	}
	tempFile.Close()

	cmd := exec.Command(editor, tempFileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		return err
	}

	err = notebook.StoreNote(note, string(content))
	if err != nil {
		return err
	}

	return nil
}

func removeNote(ctx *cli.Context) error {
	note := ctx.Args().First()
	return notebook.DeleteNote(note)
}

func listNotes(ctx *cli.Context) error {
	notesInfo, err := notebook.GetAllNoteInfo()
	if err != nil {
		return err
	}

	timeSort := ctx.Bool("time")
	reverse := ctx.Bool("reverse")

	var toSort sort.Interface
	if timeSort {
		toSort = SortByLastModified(notesInfo)
	} else {
		toSort = SortByName(notesInfo)
	}

	if reverse {
		toSort = sort.Reverse(toSort)
	}

	sort.Sort(toSort)

	for _, noteInfo := range notesInfo {
		fmt.Println(noteInfo.Name)
	}

	return nil
}

type SortByName []notebook.NoteInfo

func (f SortByName) Len() int           { return len(f) }
func (f SortByName) Less(i, j int) bool { return strings.Compare(f[i].Name, f[j].Name) < 0 }
func (f SortByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type SortByLastModified []notebook.NoteInfo

func (f SortByLastModified) Len() int           { return len(f) }
func (f SortByLastModified) Less(i, j int) bool { return f[i].LastModified.After(f[j].LastModified) } // this is backward to make latest first in lists
func (f SortByLastModified) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
