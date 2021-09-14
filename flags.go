package main

import "github.com/urfave/cli"

var (
	FlagRpcUrl *string
	FlagRawTx  *string
)

var (
	Verbose = cli.BoolFlag{
		Name:  "verbose",
		Usage: "output debug information",
	}
)
