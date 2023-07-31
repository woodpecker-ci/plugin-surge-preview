FROM --platform=$BUILDPLATFORM golang:1.20 AS build

WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-s -w -extldflags "-static"' -o plugin-surge-preview

FROM --platform=$BUILDPLATFORM node:20-alpine

RUN npm install -g surge@0.23.1
COPY --from=build src/plugin-surge-preview /bin/

ENTRYPOINT ["/bin/plugin-surge-preview"]
