FROM golang:1.12 as builder
WORKDIR /go/src/github.com/sthlmio/preemptible-sentinel
COPY vendor/ ./vendor/
COPY pkg/ ./pkg/
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o controller .

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/sthlmio/preemptible-sentinel/controller .
CMD ["/controller"]