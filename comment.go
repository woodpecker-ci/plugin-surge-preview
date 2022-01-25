package main

import (
	"context"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

const commentKey = "<!--surge.sh-->"

type comment struct {
	client *scm.Client
}

func (c *comment) Load(driver, serverUrl, oauth string) error {
	client, err := factory.NewClient(driver, serverUrl, oauth)
	if err != nil {
		return err
	}

	c.client = client
	return nil
}

func (c *comment) UpdateOrCreateComment(ctx context.Context, repo string, prID int, body string) error {
	commentInput := &scm.CommentInput{
		Body: body,
	}

	listOptions := scm.ListOptions{
		Page: 1,
		Size: 50,
	}
	for {
		comments, resp, err := c.client.PullRequests.ListComments(ctx, repo, prID, listOptions)
		if err != nil {
			return err
		}

		for _, comment := range comments {
			if strings.Contains(comment.Body, commentKey) {
				_, _, err = c.client.PullRequests.EditComment(ctx, repo, prID, comment.ID, commentInput)
				return err
			}
		}

		// Exit the loop when we've seen all pages
		if resp.Rate.Remaining == 0 {
			break
		}

		// Update the page number to get the next page
		listOptions.Page = resp.Page.Next
	}

	_, _, err := c.client.PullRequests.CreateComment(ctx, repo, prID, commentInput)
	return err
}
