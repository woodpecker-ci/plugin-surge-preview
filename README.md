# plugin-surge-preview

[![Build status](https://ci.woodpecker-ci.org/api/badges/woodpecker-ci/plugin-surge-preview/status.svg)](https://ci.woodpecker-ci.org/woodpecker-ci/plugin-surge-preview)
[![Docker Image Version (latest by date)](https://img.shields.io/docker/v/woodpeckerci/plugin-surge-preview?label=DockerHub%20latest%20version&sort=semver)](https://hub.docker.com/r/woodpeckerci/plugin-surge-preview/tags)

Woodpecker plugin to deploy static websites (for PR previews) to [surge.sh](https://surge.sh/).

## Build

Build the Docker image with the following command:

```sh
docker buildx build -t woodpeckerci/plugin-surge-preview . --load
```

## Usage

Create a Woodpecker pipeline in your project and add a step like this one:

```yml
steps:
  - name: PR preview
    image: woodpeckerci/plugin-surge-preview:<tag>
    settings:
      path: dist/ # path to directory to publish files from
      surge_token:
        from_secret: SURGE_TOKEN # install surge cli and run `surge token`: https://surge.sh/help/getting-started-with-surge
      forge_type: github # or gitea, gitlab, ... (gitea = forgejo)
      forge_url: https://github.com # or https://codeberg.org, https://gitlab.com, ...
      forge_repo_token:
        from_secret: FORGE_TOKEN # access token for your forge
    when:
      event:
        - pull_request
        - pull_request_closed
```

## Running from the CLI

```bash
docker run --rm -it \
  -e PLUGIN_PATH="dist/" \
  -e PLUGIN_SURGE_TOKEN="SURGE_TOKEN" \
  -e PLUGIN_FORGE_TYPE="gitea" \
  -e PLUGIN_FORGE_URL="https://codeberg.org" \
  -e PLUGIN_FORGE_REPO_TOKEN="FORGE_TOKEN" \
  -e CI_PIPELINE_EVENT=pull_request \
  -e CI_REPO_OWNER=REPO_OWNER \
  -e CI_REPO_NAME=REPO_NAME \
  -e CI_COMMIT_PULL_REQUEST=99 \
  -v $(pwd):/woodpecker/src \
  -w /woodpecker/src \
  --entrypoint sh \
  woodpeckerci/plugin-surge-preview:<tag>

# Inside the container, trigger the app manually:
plugin-surge-preview
# Or check the version of surge installed:
surge --version
```

To [tear down a project on surge.sh](https://surge.sh/help/tearing-down-a-project), run the CLI with the environment variable `CI_PIPELINE_EVENT=pull_request_closed`.

## Credits

- <https://github.com/afc163/surge-preview>
