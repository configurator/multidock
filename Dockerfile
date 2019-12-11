FROM golang:latest AS builder
WORKDIR /src

COPY ./go.mod ./go.sum ./*.go ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./pkg ./pkg
ENV CGO_ENABLED=0
RUN go generate ./...
RUN go build -o /multidock ./cmd/multidock

FROM scratch
COPY --from=builder /multidock /multidock

ENTRYPOINT ["/multidock"]
CMD []
