package main

import (
	"bufio"
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
	fmt.Print(game.ButtonLegend())
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	moves := 0

	for {

		//fmt.Printf("Enter text: ")
		sugg := game.Search()
		fmt.Printf("Enter text: (suggest=%v): ", sugg)
		text, _ := reader.ReadString('\n')
		b, err := flood.InputToButton(text, sugg)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			continue
		}

		fmt.Printf("Pressing button %v (%v)\n", flood.Buttons(b), b)
		moves++
		count, err := game.At.ButtonPress(flood.Buttons(b))
		fmt.Printf("count=%d   error=%v\n", count, err)
		if err != nil {
			continue
		}

		fmt.Printf("%s", game.String())
		fmt.Print(game.ButtonLegend())

		if game.Won() {
			fmt.Println("You won!!!")
			break
		}
		if moves > game.MaxMoves {
			fmt.Println("No more moves!")
			break
		}
	}

	fmt.Println()

	return nil
}
