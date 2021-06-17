package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"system-survey/cmd"
)

func main() {
	jsonFlag := &cli.BoolFlag{
		Name:        "json",
		Usage:       "Use this flag to get machine readable JSON format",
		DefaultText: "Will output in human readable format",
	}

	app := &cli.App{
		Name:        "System survey tool",
		Usage:       "",
		Description: "",
		Commands: cli.Commands{
			&cli.Command{
				Name:        "apps",
				Usage:       "Check list of installed applications",
				UsageText:   "apps [command options]",
				Description: "Reads out the PATH variable and prints all installed applications.",
				Flags: cli.FlagsByName{
					jsonFlag,
					&cli.StringFlag{
						Name:        "filter",
						Aliases:     []string{"f"},
						Usage:       "Comma separated application names",
						DefaultText: "All applications",
					},
				},
				Action: func(c *cli.Context) error {
					apps := cmd.GetAppList(c.String("filter"))
					if c.Bool("json") {
						fmt.Println(string(apps.JSON()))
					} else {
						fmt.Println(apps.String())
					}
					return nil
				},
			},
			&cli.Command{
				Name:        "golang",
				Usage:       "Check golang version",
				UsageText:   "golang [command options]",
				Description: "Runs the `go version` command and get output as a string or error on failure.",
				Flags: cli.FlagsByName{
					jsonFlag,
				},
				Action: func(c *cli.Context) error {
					gv, err := cmd.GetGolangVersion()
					if err != nil {
						return err
					}
					if c.Bool("json") {
						fmt.Println(string(gv.JSON()))
					} else {
						fmt.Println(gv.String())
					}
					return nil
				},
			},
			&cli.Command{
				Name:        "hwinfo",
				Usage:       "Get information about the host machine",
				UsageText:   "hwinfo [command options]",
				Description: "Collects system and hardware information.",
				Flags: cli.FlagsByName{
					jsonFlag,
				},
				Action: func(c *cli.Context) error {
					hi := cmd.GetHwinfo()
					if c.Bool("json") {
						fmt.Println(string(hi.JSON()))
					} else {
						fmt.Println(hi.String())
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
