package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "next"

func main() {
	app := &cli.Command{}
	app.Name = "surge-preview plugin"
	app.Usage = "surge-preview plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "surge-token",
			Usage:   "token for surge authentication",
			Sources: cli.EnvVars("PLUGIN_SURGE_TOKEN"),
		},
		&cli.StringFlag{
			Name:    "path",
			Usage:   "path to upload files from",
			Sources: cli.EnvVars("PLUGIN_PATH"),
		},
		&cli.StringFlag{
			Name:    "pipeline-event",
			Usage:   "event of the current pipeline",
			Sources: cli.EnvVars("CI_PIPELINE_EVENT", "CI_BUILD_EVENT"),
		},
		&cli.StringFlag{
			Name:    "repo-owner",
			Usage:   "owner of the current repo",
			Sources: cli.EnvVars("CI_REPO_OWNER", "DRONE_REPO_OWNER"),
		},
		&cli.StringFlag{
			Name:    "repo-name",
			Usage:   "name of the current repo",
			Sources: cli.EnvVars("CI_REPO_NAME", "DRONE_REPO_NAME"),
		},
		&cli.StringFlag{
			Name:    "pull-request-id",
			Usage:   "id of the current pull-request",
			Sources: cli.EnvVars("CI_COMMIT_PULL_REQUEST", "DRONE_PULL_REQUEST"),
		},
		&cli.StringFlag{
			Name:    "forge-type",
			Usage:   "type of the forge",
			Sources: cli.EnvVars("CI_FORGE_TYPE", "PLUGIN_FORGE_TYPE"),
		},
		&cli.StringFlag{
			Name:    "forge-url",
			Usage:   "base url of the forge",
			Sources: cli.EnvVars("CI_FORGE_URL", "PLUGIN_FORGE_URL"),
		},
		&cli.StringFlag{
			Name:    "forge-repo-token",
			Usage:   "type of the forge",
			Sources: cli.EnvVars("CI_FORGE_REPO_TOKEN", "PLUGIN_FORGE_REPO_TOKEN"),
		},
	}

	ctx := context.Background()
	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, c *cli.Command) error {
	plugin := Plugin{
		SurgeToken:     c.String("surge-token"),
		Path:           c.String("path"),
		RepoOwner:      c.String("repo-owner"),
		RepoName:       c.String("repo-name"),
		PipelineEvent:  c.String("pipeline-event"),
		PullRequestID:  int(c.Int("pull-request-id")),
		ForgeType:      c.String("forge-type"),
		ForgeURL:       c.String("forge-url"),
		ForgeRepoToken: c.String("forge-repo-token"),
	}

	return plugin.Exec(ctx)
}
