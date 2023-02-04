package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func enablecli() {
	appCommands := []*cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "open config file",
			Action: func(ctx *cli.Context) error {
				// open .yaml with default editor, otherwise just print path
				if strings.HasSuffix(configpath, ".yaml") {
					open.Run(configpath)
				} else {
					fmt.Printf("config file location: %v\n", configpath)
				}
				return nil
			},
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "set a secret `key=value`",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					pair := c.Args().First()
					key, value, foundEqual := strings.Cut(pair, "=")

					// check if secret already exists and prompt for overwrite
					if secretService.HasSecretKey(key) {
						fmt.Print("Overwrite value? (y/n) ")
						var choice string
						fmt.Scanln(&choice)
						if choice != "y" && choice != "Y" {
							return nil
						}
					}

					// Read secret value since argument was key only
					if !foundEqual {
						fmt.Printf("Enter value: ")
						fd := os.Stdin.Fd()
						bRead, err := term.ReadPassword(int(fd))
						if err != nil {
							log.Fatal(err)
						}
						value = string(bRead)
						fmt.Println()
					}

					// delete secret if it already exists
					if secretService.HasSecretKey(key) {
						err := secretService.DeleteSecret(key)
						if err != nil {
							return err
						}
					}
					err := secretService.AddSecret(Secret{key, value})
					if err != nil {
						return err
					}
					log.Infof("%v saved", key)
				} else {
					return errors.New("`key=value` pair not found to save")
				}
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get secret for a `key`",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					key := c.Args().First()
					value, err := secretService.GetSecretValue(key)
					if err != nil {
						return err
					}
					if c.NArg() > 1 {
						err = execCmd(c.Args().Slice()[1], c.Args().Slice()[2:], []string{key + "=" + value})

						if err != nil {
							return err
						}
					} else {
						fmt.Println(value)
					}
				} else {
					return errors.New("`key=value` pair not found to save")
				}
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all secrets saved using `uses`",
			Action: func(ctx *cli.Context) error {
				keys := secretService.ListSecretKeys()
				fmt.Println(strings.Join(keys, "\n"))
				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "delete a `secret`",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() > 0 {
					key := ctx.Args().First()
					err := secretService.DeleteSecret(key)
					if err != nil {
						return err
					}
				} else {
					return errors.New("`key` required for removal")
				}
				return nil
			},
		},
		{
			Name:  "completion",
			Usage: "generate auto-complete commands for a shell",
			Subcommands: []*cli.Command{
				{
					Name: "zsh",
					Action: func(ctx *cli.Context) error {
						return generateAutocomplete("zsh")
					},
				},
				{
					Name: "bash",
					Action: func(ctx *cli.Context) error {
						return generateAutocomplete("bash")
					},
				},
			},
		},
	}
	for _, m := range config.Project {
		projectCommands := []*cli.Command{{
			Name:  m.Name,
			Usage: "get secrets for project `" + m.Name + "` and run command",
			Action: func(ctx *cli.Context) error {
				// find project with name
				getProjectSecrets := func(name string) ([]ConfigSecret, error) {
					for _, p := range config.Project {
						if p.Name == name {
							return p.ConfigSecret[:], nil
						}
					}
					return nil, errors.New("unable to find in config project " + name)
				}
				secretNames, err := getProjectSecrets(ctx.Command.Name)
				if err != nil {
					log.Fatal(err)
				}

				// read all secrets listed in project
				count := len(secretNames)
				secrets := make([]string, count)
				var value string
				for i, configSecret := range secretNames {
					value, err = secretService.GetSecretValue(configSecret.Key)
					if err != nil {
						log.Errorf("secret not found %v", configSecret.Key)
						log.Fatal(err)
					}
					secrets[i] = configSecret.VariableName + "=" + value
				}

				// run command with secrets
				if ctx.NArg() > 0 {
					err = execCmd(ctx.Args().Slice()[0], ctx.Args().Slice()[1:], secrets)
					if err != nil {
						return err
					}
				} else {
					return errors.New("no command to run using project secrets")
				}

				return nil
			},
		},
		}

		// ensure that project commands show up first
		appCommands = append(projectCommands, appCommands...)
	}
	app := &cli.App{
		Name:                 "uses",
		Usage:                "securely manage secrets in dev environment",
		Commands:             appCommands,
		Version:              Version,
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
