package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
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
			Action:  listNotes,
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

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't edit note '%s' with %s", note, editor)
	}
	return err
}

func removeNote(ctx *cli.Context) error {
	fmt.Println("remove " + ctx.Args().First())
	return nil
}

func listNotes(ctx *cli.Context) error {
	fmt.Println("list notes")
	return nil
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
