package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaanek/jeth/ui"
	"github.com/urfave/cli"
)

type Command func(term ui.Screen, ctx *cli.Context) error

var (
	app = NewApp("rarity bot")
)

func init() {
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "runs a bot",
			Action: runCommand(RunCommand),
			Flags: []cli.Flag{
				Verbose,
			},
		},
	}
}

func runCommand(cmd Command) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		term := ui.NewTerminal(false)
		err := cmd(term, ctx)
		if err != nil {
			term.Error(err)
		}
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		code := 1
		fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
}

// NewApp creates an app with sane defaults.
func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	app.Email = ""
	app.Usage = usage
	return app
}
