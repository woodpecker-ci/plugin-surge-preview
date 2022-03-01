package main

import (
	"context"
	"errors"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

const commentKey = "<!--woodpeckerci-plugin-surge-preview-->"

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

func (c *comment) Find(ctx context.Context, repo string, prID int) (*scm.Comment, error) {
	listOptions := scm.ListOptions{
		Page: 1,
		Size: 1,
	}
	for {
		comments, resp, err := c.client.PullRequests.ListComments(ctx, repo, prID, listOptions)
		if err != nil {
			return nil, err
		}

		for _, comment := range comments {
			if strings.Contains(comment.Body, commentKey) {
				return comment, nil
			}
		}

		// Exit the loop when we've seen all pages
		if resp.Page.Next == 0 {
			break
		}

		if resp.Rate.Remaining == 0 {
			return nil, errors.New("Rate limit exceeded")
		}

		// Update the page number to get the next page
		listOptions.Page = resp.Page.Next
	}

	return nil, errors.New("Comment not found")
}

func (c *comment) UpdateOrCreateComment(ctx context.Context, repo string, prID int, comment *scm.Comment, body string) (*scm.Comment, error) {
	commentInput := &scm.CommentInput{
		Body: body,
	}

	if comment == nil {
		comment, _, err := c.client.PullRequests.CreateComment(ctx, repo, prID, commentInput)
		return comment, err
	}

	comment, _, err := c.client.PullRequests.EditComment(ctx, repo, prID, comment.ID, commentInput)
	return comment, err
}
