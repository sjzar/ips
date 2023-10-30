FROM golang:1.21 as builder

COPY . /building
WORKDIR /building
RUN make build

FROM alpine:3 as alpine

WORKDIR /
COPY --from=builder /building/bin/ips /bin/ips
EXPOSE 6860
ENTRYPOINT ["/bin/ips", "server"]
