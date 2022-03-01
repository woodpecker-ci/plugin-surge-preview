package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Plugin struct {
	SurgeToken     string
	Path           string
	RepoOwner      string
	RepoName       string
	PipelineEvent  string
	PullRequestId  int
	ForgeType      string
	ForgeUrl       string
	ForgeRepoToken string
	comment        *comment
}

func (p *Plugin) Exec() error {
	fmt.Println("Surge.sh preview plugin")

	if p.RepoName == "" || p.RepoOwner == "" || p.PipelineEvent == "" {
		return errors.New("Missing required parameters. Are you running this plugin from within a pipeline?")
	}

	if p.Path == "" {
		return errors.New("Path to the upload folder is not set")
	}

	if p.SurgeToken == "" {
		return errors.New("Surge.sh token is not defined")
	}

	p.comment = &comment{}
	p.comment.Load(p.ForgeType, p.ForgeUrl, p.ForgeRepoToken)

	switch p.PipelineEvent {
	case "pull_request":
		return p.deploy()
	case "pull_close":
		return p.teardown()
	default:
		return errors.New("unsupported pipeline event, please only run on pull_request or pull_close")
	}
}

func (p *Plugin) deploy() error {
	url := p.getPreviewUrl()
	repo := p.RepoOwner + "/" + p.RepoName
	fmt.Printf("Deploying preview to %s\n", url)
	p.comment.UpdateOrCreateComment(context.Background(), repo, p.PullRequestId, "Deploying preview to: "+url)

	if err := p.runSurgeCommand(false); err != nil {
		return err
	}

	fmt.Println("Deployment of preview was successful")
	p.comment.UpdateOrCreateComment(context.Background(), repo, p.PullRequestId, "Deployment of preview was successful: "+url)

	return nil
}

func (p *Plugin) teardown() error {
	url := p.getPreviewUrl()
	repo := p.RepoOwner + "/" + p.RepoName
	fmt.Printf("Teading down %s", url)

	if err := p.runSurgeCommand(true); err != nil {
		return err
	}

	fmt.Println("Preview torn down")
	p.comment.UpdateOrCreateComment(context.Background(), repo, p.PullRequestId, "Deployment of preview was torn down")

	return nil
}

func (p *Plugin) getPreviewUrl() string {
	return fmt.Sprintf("%s-%s-pr-%d.surge.sh", p.RepoOwner, p.RepoName, p.PullRequestId)
}

func (p *Plugin) runSurgeCommand(teardown bool) error {
	cmdArg := "./" + p.Path

	if teardown {
		cmdArg = "teardown"
	}

	cmd := exec.Command("surge", cmdArg, p.getPreviewUrl(), `--token`, p.SurgeToken)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(stdout))

	if !strings.Contains(string(stdout), "Success!") {
		return errors.New("Failed to run surge")
	}

	return nil
}
