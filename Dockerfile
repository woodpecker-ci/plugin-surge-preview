FROM --platform=$BUILDPLATFORM golang:1.21@sha256:24a09375a6216764a3eda6a25490a88ac178b5fcb9511d59d0da5ebf9e496474 AS build

WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-s -w -extldflags "-static"' -o plugin-surge-preview

FROM --platform=$BUILDPLATFORM node:21-alpine@sha256:df76a9449df49785f89d517764012e3396b063ba3e746e8d88f36e9f332b1864

RUN npm install -g surge@0.23.1
COPY --from=build src/plugin-surge-preview /bin/

ENTRYPOINT ["/bin/plugin-surge-preview"]
