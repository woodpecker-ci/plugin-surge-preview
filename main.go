package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "next"
)

func main() {
	app := cli.NewApp()
	app.Name = "surge-preview plugin"
	app.Usage = "surge-preview plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "surge-token",
			Usage:   "token for surge authentication",
			EnvVars: []string{"PLUGIN_SURGE_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "path",
			Usage:   "path to upload files from",
			EnvVars: []string{"PLUGIN_PATH"},
		},
		&cli.StringFlag{
			Name:    "pipeline-event",
			Usage:   "event of the current pipeline",
			EnvVars: []string{"CI_BUILD_EVENT"},
		},
		&cli.StringFlag{
			Name:    "repo-owner",
			Usage:   "owner of the current repo",
			EnvVars: []string{"CI_REPO_OWNER"},
		},
		&cli.StringFlag{
			Name:    "repo-name",
			Usage:   "name of the current repo",
			EnvVars: []string{"CI_REPO_NAME"},
		},
		&cli.StringFlag{
			Name:    "pull-request-id",
			Usage:   "id of the current pull-request",
			EnvVars: []string{"CI_COMMIT_PULL_REQUEST"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		SurgeToken:    c.String("surge-token"),
		Path:          c.String("path"),
		RepoOwner:     c.String("repo-owner"),
		RepoName:      c.String("repo-name"),
		PipelineEvent: c.String("pipeline-event"),
		PipelineLink:  c.String("pipeline-link"),
		PullRequestId: c.Int("pull-request-id"),
	}

	return plugin.Exec()
}
