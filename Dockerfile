# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder

RUN apt-get update
RUN apt-get install -y less vim

RUN useradd -ms /bin/bash testground
USER testground
WORKDIR /home/testground

# Copy local code to the container image.
COPY --chown=testground Makefile go.mod go.sum main.go ./
COPY --chown=testground static static
RUN make checkout build-targets

WORKDIR /home/testground/testground/plans/smlbench/cmd
RUN go build -v -o smlbench
WORKDIR /home/testground

ENV PATH="/home/testground/targets/go-ipfs/cmd/ipfs:${PATH}"

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go build -v -o testground-web

# Run the web service on container startup.
EXPOSE 8011
CMD ["/home/testground/testground-web"]
