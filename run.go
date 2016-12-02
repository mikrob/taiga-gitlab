package main

import (
	"os"

	"taiga-gitlab/importers"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "taiga-gitlab"
	app.Description = "Taiga <-> Gitlab"
	app.Version = "0.1.0"
	cliGitlab := []cli.Flag{
		cli.StringFlag{
			Name:  "gitlab-url",
			Usage: "Gitlab URL",
		},
		cli.StringFlag{
			Name:   "gitlab-token",
			Usage:  "Gitlab API Token",
			EnvVar: "GITLAB_TOKEN",
		},
		cli.StringFlag{
			Name:  "gitlab-project",
			Usage: "Gitlab project name (x/y form)",
		},
	}
	cliTaiga := []cli.Flag{
		cli.StringFlag{
			Name:  "taiga-url",
			Usage: "Taiga URL",
		},
		cli.StringFlag{
			Name:   "taiga-user",
			Usage:  "Taiga username",
			EnvVar: "TAIGA_USER",
		},
		cli.StringFlag{
			Name:   "taiga-password",
			Usage:  "Taiga password",
			EnvVar: "TAIGA_PASSWORD",
		},
		cli.StringFlag{
			Name:  "taiga-project",
			Usage: "Taiga project name",
		},
		cli.BoolFlag{
			Name:  "taiga-skip-user",
			Usage: "Skip user creation",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "import",
			Usage:  "import Gitlab issues to Taiga",
			Action: importers.ImportGitlab2Taiga,
			Flags:  append(cliGitlab, cliTaiga...),
		},
	}
	app.Run(os.Args)
}
