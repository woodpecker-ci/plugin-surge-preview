# plugin-surge-preview

Woodpecker plugin to deploy static pages for reviewing to [surge.sh](https://surge.sh/).

## Build

Build the binary with the following command:

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -ldflags '-s -w -extldflags "-static"' -o plugin-codecov
```

## Docker

Build the Docker image with the following command:

```sh
docker build -f docker/Dockerfile -t woodpeckerci/plugin-surge-preview .
```
