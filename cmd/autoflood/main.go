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
	app.Name = "autoflood"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{
			Name:    "stats",
			Aliases: []string{"s"},
			Usage:   "print stats about this game board",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return errors.Errorf("Must specify 1 file on the command line.")
				}
				filename := c.Args().Get(0)
				fmt.Printf("FILE=%q\n", filename)

				_, err := initialize(filename)
				return err
			},
		},
		{
			Name:    "play1",
			Aliases: []string{"1", "play"},
			Usage:   "play the game with simple suggestions",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return errors.Errorf("Must specify 1 file on the command line.")
				}
				filename := c.Args().Get(0)
				fmt.Printf("FILE=%q\n", filename)

				game, err := initialize(filename)
				if err != nil {
					return err
				}

				return playManual(game, game.Search1)
			},
		},
		{
			Name:    "play2",
			Aliases: []string{"2"},
			Usage:   "play the game with multi-level suggestions",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return errors.Errorf("Must specify 1 file on the command line.")
				}
				filename := c.Args().Get(0)
				fmt.Printf("FILE=%q\n", filename)

				game, err := initialize(filename)
				if err != nil {
					return err
				}

				return playManual(game, game.SearchMultiLevel)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func initialize(filename string) (game flood.Game, err error) {

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

	return
}

func playManual(game flood.Game, suggest func() flood.Buttons) (err error) {

	reader := bufio.NewReader(os.Stdin)

	moves := 0

	for {

		//fmt.Printf("Enter text: ")
		sugg := suggest()
		fmt.Printf("%d) Enter text: (suggest=%v): ", moves, sugg)
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
