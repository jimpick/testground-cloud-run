# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/jimpick/testground-cloud-run
COPY . .
COPY static static

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o testground-web

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/jimpick/testground-cloud-run/testground-web /testground-cloud-run/testground-web
COPY --from=builder /go/src/github.com/jimpick/testground-cloud-run/static /testground-cloud-run/static

# Run the web service on container startup.
WORKDIR /testground-cloud-run
EXPOSE 8011
CMD ["/testground-cloud-run/testground-web"]
