package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gowallet"
	app.Commands = []cli.Command{
		{
			Name:    "transfer",
			Aliases: []string{"t"},
			Usage:   "transfer wei to address",
			Action: func(c *cli.Context) error {
				fmt.Println("transfer")
				return nil
			},
		},
		{
			Name:    "balance",
			Aliases: []string{"b"},
			Usage:   "check wallet balance",
			Action: func(c *cli.Context) error {
				fmt.Println("balance")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
