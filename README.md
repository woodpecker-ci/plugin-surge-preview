# plugin-surge-preview

Woodpecker plugin to deploy static pages for reviewing to [surge.sh](https://surge.sh/).

## Build

Build the Docker image with the following command:

```sh
docker buildx build -t woodpeckerci/plugin-surge-preview . --load
```

## Usage

Create a woodpecker pipeline in your project and add a step like this one:

```yml
pipeline:
  preview:
    image: woodpeckerci/plugin-surge-preview
    settings:
      path: dist/ # path to directory to publish files from
      surge_token:
        from_secret: SURGE_TOKEN # install surge cli and run `surge token`: https://surge.sh/help/getting-started-with-surge
      forge_type: github # or gitea, gitlab, ...
      forge_url: https://github.com # or https://codeberg.org, https://gitlab.com, ...
      forge_repo_token:
        from_secret: FORGE_TOKEN # access token for your forge
    when:
      event: pull_request
```

## Running from the CLI

```bash
docker run --rm -it \
  -e PLUGIN_PATH="dist/" \
  -e PLUGIN_SURGE_TOKEN="SURGE_TOKEN" \
  -e PLUGIN_FORGE_TYPE="gitea" \
  -e PLUGIN_FORGE_URL="https://codeberg.org" \
  -e PLUGIN_FORGE_REPO_TOKEN="FORGE_TOKEN" \
  -e CI_BUILD_EVENT=pull_request \
  -e CI_REPO_OWNER=REPO_OWNER \
  -e CI_REPO_NAME=REPO_NAME \
  -e CI_COMMIT_PULL_REQUEST=99 \
  -v $(pwd):/drone/src \
  -w /drone/src \
  --entrypoint sh \
  woodpeckerci/plugin-surge-preview:next

# Inside the container, trigger the app manually:
plugin-surge-preview
# Or check the version of surge installed:
surge --version
```

To [tear down a project on surge.sh](https://surge.sh/help/tearing-down-a-project), run the CLI with the environment variable `CI_BUILD_EVENT=pull_close`.

## Credits

- https://github.com/afc163/surge-preview
