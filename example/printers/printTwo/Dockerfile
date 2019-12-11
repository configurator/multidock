FROM golang:latest AS builder
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build -a -o /print ./print.go

FROM scratch
COPY --from=builder /print /print
CMD ["/print"]
