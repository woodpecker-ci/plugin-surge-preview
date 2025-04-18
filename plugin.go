package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Plugin struct {
	SurgeToken     string
	Path           string
	RepoOwner      string
	RepoName       string
	PipelineEvent  string
	PullRequestID  int
	ForgeType      string
	ForgeURL       string
	ForgeRepoToken string
	comment        *comment
}

func (p *Plugin) Exec(ctx context.Context) error {
	fmt.Println("--- Surge.sh preview plugin ---")
	fmt.Println()

	if p.RepoName == "" || p.RepoOwner == "" || p.PipelineEvent == "" {
		return errors.New("Missing required parameters. Are you running this plugin from within a pipeline?")
	}

	if p.Path == "" {
		return errors.New("Path to the upload folder is not set")
	}

	if p.SurgeToken == "" {
		return errors.New("Surge.sh token is missing")
	}

	p.comment = &comment{}

	driver := p.ForgeType
	// https://github.com/jenkins-x/go-scm does not understand forgejo.
	if driver == "forgejo" {
		driver = "gitea"
	}
	err := p.comment.Load(driver, p.ForgeURL, p.ForgeRepoToken)
	if err != nil {
		return err
	}

	switch p.PipelineEvent {
	case "pull_request":
		return p.deploy(ctx)
	case "pull_request_closed":
		return p.teardown(ctx)
	default:
		return errors.New("unsupported pipeline event, only 'pull_request' and 'pull_request_closed' are supported")
	}
}

func (p *Plugin) deploy(ctx context.Context) error {
	url := p.getPreviewURL()
	repo := p.RepoOwner + "/" + p.RepoName

	comment, err := p.comment.Find(ctx, repo, p.PullRequestID)
	if err != nil && err.Error() != "Comment not found" {
		return err
	}

	commentText := fmt.Sprintf("Deploying preview to https://%s", url)
	fmt.Println(commentText)
	comment, err = p.comment.UpdateOrCreateComment(ctx, repo, p.PullRequestID, comment, commentText)
	if err != nil {
		return err
	}

	if err := p.runSurgeCommand(false); err != nil {
		return err
	}

	commentText = fmt.Sprintf("Surge PR preview deployment succeeded. View it at https://%s", url)
	fmt.Println(commentText)
	_, err = p.comment.UpdateOrCreateComment(ctx, repo, p.PullRequestID, comment, commentText)
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) teardown(ctx context.Context) error {
	url := p.getPreviewURL()
	repo := p.RepoOwner + "/" + p.RepoName

	comment, err := p.comment.Find(ctx, repo, p.PullRequestID)
	if err != nil && err.Error() != "Comment not found" {
		return err
	}

	commentText := fmt.Sprintf("Shutting down https://%s\n", url)
	fmt.Println(commentText)
	comment, err = p.comment.UpdateOrCreateComment(ctx, repo, p.PullRequestID, comment, commentText)
	if err != nil {
		return err
	}

	if err := p.runSurgeCommand(true); err != nil {
		return err
	}

	commentText = "Surge PR preview deployment was removed"
	fmt.Println(commentText)
	_, err = p.comment.UpdateOrCreateComment(ctx, repo, p.PullRequestID, comment, commentText)
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) getPreviewURL() string {
	pattern := regexp.MustCompile(`[^a-zA-Z0-9\-]`)
	owner := pattern.ReplaceAllString(p.RepoOwner, "-")
	repo := pattern.ReplaceAllString(p.RepoName, "-")
	return fmt.Sprintf("%s-%s-pr-%d.surge.sh", owner, repo, p.PullRequestID)
}

func (p *Plugin) runSurgeCommand(teardown bool) error {
	var output bytes.Buffer
	var waitGroup sync.WaitGroup

	cmdArg := p.Path

	if teardown {
		cmdArg = "teardown"
	}

	cmd := exec.Command("surge", cmdArg, p.getPreviewURL(), `--token`, p.SurgeToken)
	fmt.Println("#", strings.Join(append(cmd.Args, p.getPreviewURL(), "--token ****"), " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	writer := io.MultiWriter(os.Stdout, &output)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		_, _ = io.Copy(writer, stdout)
	}()

	if err := cmd.Run(); err != nil {
		return err
	}
	waitGroup.Wait()

	if !strings.Contains(output.String(), "Success") {
		return errors.New("Failed to run surge")
	}

	return nil
}
