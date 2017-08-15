package main

import (
	"fmt"
	"os"

	"github.com/TomOnTime/autoflood/pkg/extractflood"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			return errors.Errorf("Must specify 1 file on the command line.")
		}
		filename := c.Args().Get(0)
		fmt.Printf("FILE=%q\n", filename)
		return extractflood.ExtractFile(filename)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
