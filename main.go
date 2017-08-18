package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "envdo"
	app.Usage = "Manage environment variables efficiently."
	app.Version = "0.1.1"
	app.Author = "Anton Lindstrom"

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "homedir: %s\n", err)
		os.Exit(2)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "directory,d",
			Value:  homeDir + "/.envdo",
			Usage:  "path to directory with environment variables",
			EnvVar: "ENVDO_DIR",
		},
		cli.StringFlag{
			Name:   "gpg-recipient,r",
			Usage:  "GPG recipient to encrypt or decrypt environment variables (email or ID)",
			EnvVar: "ENVDO_GPG_RECIPIENT",
		},
		cli.BoolFlag{
			Name:  "preserve-env,E",
			Usage: "preserve user environment when running command",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}

		if c.NArg() < 2 {
			fmt.Printf("usage: %s <profile> <command>\n", os.Args[0])
			fmt.Println("\nExample usage:")
			fmt.Println("  envdo Business/cheese-whiz env # Run env with profile Business/cheese-whiz.")
			fmt.Println("  envdo Home/backup ls -1 -h     # Run ls -1 -h with profile Home/backup.")
			fmt.Println("")
			return nil
		}

		var gpgRecipient *string
		if c.GlobalIsSet("gpg-recipient") {
			gpgRecipient = func(s string) *string { return &s }(c.GlobalString("gpg-recipient"))
		}

		profileEnvVars, err := fetchEnv(c.GlobalString("directory"), c.Args()[0], gpgRecipient)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		err = executeCmd(c.Args()[1], c.Args()[2:], c.GlobalBool("preserve-env"), profileEnvVars)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	app.Commands = []cli.Command{
		addCmd,
		lsCmd,
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

var lsCmd = cli.Command{
	Name:  "ls",
	Usage: "List all profiles",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "plain",
			Usage: "plain text, without colors or indent",
		},
	},
	Action: func(c *cli.Context) error {
		var extraPath string
		if len(c.Args()) > 0 {
			extraPath = "/" + c.Args()[0]
		}

		if c.Bool("plain") {
			files, err := readFiles(c.GlobalString("directory") + extraPath)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			for _, f := range files {
				fmt.Println(profileName(c.GlobalString("directory")+extraPath, f))
			}
			return nil
		}

		printTree(c.GlobalString("directory") + extraPath)

		return nil
	},
}

var addCmd = cli.Command{
	Name:  "add",
	Usage: "add a GPG2 encrypted environment variable file",
	Action: func(c *cli.Context) error {
		if c.GlobalString("gpg-recipient") == "" {
			return cli.NewExitError("error: gpg recipient is required (-r)", 1)
		}

		profileName := c.Args()[0] + ".gpg"
		absoluteProfile := fmt.Sprintf("%s/%s", c.GlobalString("directory"), profileName)
		file := filepath.Base(profileName)

		directory := strings.TrimSuffix(absoluteProfile, file)
		err := os.MkdirAll(directory, 0700)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Printf("Format for environment variables is `KEY=value`.\n")
		fmt.Printf("Enter contents for %s and press Ctrl+D when finished:\n", profileName)

		gpgArgs := []string{"-e", "-r", c.GlobalString("gpg-recipient"), "-o", absoluteProfile}

		err = executeCmd("gpg2", gpgArgs, true, nil)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	},
}
