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

## Credits

- https://github.com/afc163/surge-preview
