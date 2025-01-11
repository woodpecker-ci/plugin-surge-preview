FROM --platform=$BUILDPLATFORM golang:1.23 AS build

WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-s -w -extldflags "-static"' -o plugin-surge-preview

FROM --platform=$BUILDPLATFORM node:22-alpine

# renovate: datasource=github-tags depName=sintaxi/surge
ENV SURGE_VERSION=v0.24.6

RUN npm install -g surge@${SURGE_VERSION#v}
COPY --from=build src/plugin-surge-preview /bin/

ENTRYPOINT ["/bin/plugin-surge-preview"]
