package main

import (
	"fmt"
	"os"

	"github.com/TomOnTime/autoflood/pkg/flood"
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
		return play(filename)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func play(filename string) (err error) {
	var game flood.Game

	err = game.LoadImage(filename)
	if err != nil {
		return
	}

	err = game.IdentifyLevel()
	if err != nil {
		return
	}

	err = game.ExtractGrid()
	if err != nil {
		return
	}

	err = game.ExtractButtons()
	if err != nil {
		return
	}

	fmt.Printf("%s", game.String())
	for i, n := range game.ButtonNames {
		fmt.Printf("%s:%s ", flood.Buttons(i), n)
	}
	fmt.Println()

	return nil
}
