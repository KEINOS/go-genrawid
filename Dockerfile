# =============================================================================
#  Docker Example
# =============================================================================
#  This is a Dockerfile for those who do not have Go installed locally.
#
#  - How to build:
#      docker build -t genrawid:local .
#
#  - How to run:
#      # Show help
#      docker run genrawid:local --help
#
#      # Get the rawid of a string via arg
#      docker run -i genrawid:local --string "abcdefgh"
#
#      # Get the rawid of a string via piped STDIN
#      echo -n 'abcdefgh' | docker run -i genrawid:local -
#
#      # Get the rawid of a file via piped STDIN. The below two are equivalent.
#      cat ./my_sample.txt | docker run -i genrawid:local -
#      docker run -i genrawid:local - < ./my_sample.txt
#
#  - Note:
#      To get the rawid of a file, you need to mount the file to the container.
#      Though, piped input as shown above is recommended.
# =============================================================================

# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:alpine AS build

COPY . /workspace
WORKDIR /workspace

ENV CGO_ENABLED=0

RUN go mod download
RUN go install ./cmd/genrawid

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM scratch

COPY --from=build /go/bin/genrawid /usr/bin/genrawid

ENTRYPOINT [ "/usr/bin/genrawid" ]
