---
name: Surge preview plugin
authors: Woodpecker Authors
icon: https://woodpecker-ci.org/img/logo.svg
description: Plugin to create static pages deployments as preview environments on pull-requests.
tags: [publish, cdn, preview]
containerImage: woodpeckerci/plugin-surge-preview
containerImageUrl: https://hub.docker.com/r/woodpeckerci/plugin-surge-preview
url: https://github.com/woodpecker-ci/plugin-surge-preview
---

# plugin-surge-preview

The surge-preview plugin uploads a files of a directory to the CDN of [surge.sh](https://surge.sh/) it automatically generates an url and posts the status of the deployment with an url as a comment to the pull-request. After closing a pull-request it automatically destroys the preview environment again.

## Usage

To use the plugin add a step similar to the following one to your Woodpecker pipeline config:

```yml
pipeline:
  preview:
    image: woodpeckerci/plugin-surge-preview
    settings:
      path: dist/ # path to directory to publish files from
      surge_token: xxx # install surge cli and run `surge token`: https://surge.sh/help/getting-started-with-surge
      forge_type: github # or gitea, gitlab, ...
      forge_url: https://github.com # or https://codeberg.org, https://gitlab.com, ...
      forge_repo_token: xxx # access token for your forge
    when:
      event: pull_request
```
