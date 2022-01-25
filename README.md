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
      path: dist/
      surge_token: xxx
      scm_token: xxx # access token for your SCM
      scm_type: github # or gitea, gitlab, ...
    when:
      event: pull_request
```

## Credits

- https://github.com/afc163/surge-preview